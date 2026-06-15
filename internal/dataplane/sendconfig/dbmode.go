package sendconfig

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/blang/semver/v4"
	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/diff"
	"github.com/kong/go-database-reconciler/pkg/dump"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/kong/go-database-reconciler/pkg/state"
	deckutils "github.com/kong/go-database-reconciler/pkg/utils"
	"github.com/kong/go-kong/kong"
	"github.com/samber/mo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/deckerrors"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/metrics"
)

// UpdateStrategyDBMode implements the UpdateStrategy interface. It updates Kong's data-plane
// configuration using decK's syncer.
type UpdateStrategyDBMode struct {
	client            *kong.Client
	dumpConfig        dump.Config
	version           semver.Version
	concurrency       int
	diagnostic        *diagnostics.Client
	isKonnect         bool
	logger            logr.Logger
	resourceErrors    []ResourceError
	resourceErrorLock sync.Mutex
}

// UpdateStrategyDBModeOpt is a functional option for UpdateStrategyDBMode.
type UpdateStrategyDBModeOpt func(*UpdateStrategyDBMode)

// WithDiagnostic sets the diagnostic server to send diffs to.
func WithDiagnostic(diagnostic *diagnostics.Client) UpdateStrategyDBModeOpt {
	return func(s *UpdateStrategyDBMode) {
		s.diagnostic = diagnostic
	}
}

func NewUpdateStrategyDBMode(
	client *kong.Client,
	dumpConfig dump.Config,
	version semver.Version,
	concurrency int,
	logger logr.Logger,
	opts ...UpdateStrategyDBModeOpt,
) *UpdateStrategyDBMode {
	s := &UpdateStrategyDBMode{
		client:      client,
		dumpConfig:  dumpConfig,
		version:     version,
		concurrency: concurrency,
		logger:      logger,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func NewUpdateStrategyDBModeKonnect(
	client *kong.Client,
	dumpConfig dump.Config,
	version semver.Version,
	concurrency int,
	logger logr.Logger,
) *UpdateStrategyDBMode {
	s := NewUpdateStrategyDBMode(client, dumpConfig, version, concurrency, logger)
	s.isKonnect = true
	return s
}

func (s *UpdateStrategyDBMode) Update(ctx context.Context, targetContent ContentWithHash) (mo.Option[int], error) {
	cs, err := s.currentState(ctx)
	if err != nil {
		return mo.None[int](), fmt.Errorf("failed getting current state for %s: %w", s.client.BaseRootURL(), err)
	}

	ts, err := s.targetState(ctx, cs, targetContent.Content)
	if err != nil {
		return mo.None[int](), deckerrors.ConfigConflictError{Err: err}
	}

	if err := s.refillPluginIDs(cs, ts); err != nil {
		return mo.None[int](), err
	}

	if s.isKonnect {
		if err := refillCredentialIDs(cs, ts, s.logger); err != nil {
			return mo.None[int](), err
		}
	}

	syncer, err := diff.NewSyncer(diff.SyncerOpts{
		CurrentState:        cs,
		TargetState:         ts,
		KongClient:          s.client,
		SilenceWarnings:     true,
		IsKonnect:           s.isKonnect,
		IncludeLicenses:     true,
		EnableEntityActions: true,
	})
	if err != nil {
		return mo.None[int](), fmt.Errorf("creating a new syncer for %s: %w", s.client.BaseRootURL(), err)
	}

	ctx, cancel := context.WithCancel(ctx)
	go s.HandleEvents(ctx, syncer.GetResultChan(), s.diagnostic, fmt.Sprintf("%x", targetContent.Hash))

	_, errs, _ := syncer.Solve(ctx, s.concurrency, false, false)
	cancel()
	s.resourceErrorLock.Lock()
	defer s.resourceErrorLock.Unlock()
	resourceFailures := resourceErrorsToResourceFailures(s.resourceErrors, s.logger)
	if errs != nil {
		return mo.None[int](), NewUpdateErrorWithoutResponseBody(
			resourceFailures,
			deckutils.ErrArray{Errors: errs},
		)
	}

	// as of GDR 1.8 we should always get a plain error set in addition to resourceErrors, so returning resourceErrors
	// here should not be necessary. Return it anyway as a future-proof because why not.
	if len(resourceFailures) > 0 {
		return mo.None[int](), NewUpdateErrorWithoutResponseBody(
			resourceFailures,
			errors.New("go-database-reconciler found resource errors"),
		)
	}
	// For DB-mode there is no size to return, so we return None in case of success too.
	return mo.None[int](), nil
}

// HandleEvents handles logging and error reporting for individual entity change events generated during a sync by
// looping over an event channel. It terminates when its context dies.
func (s *UpdateStrategyDBMode) HandleEvents(
	ctx context.Context,
	events chan diff.EntityAction,
	diagnostic *diagnostics.Client,
	hash string,
) {
	s.resourceErrorLock.Lock()
	diff := diagnostics.ConfigDiff{
		Hash:     hash,
		Entities: []diagnostics.EntityDiff{},
	}
	for {
		select {
		case event := <-events:
			if event.Error == nil {
				// TODO https://github.com/Kong/go-database-reconciler/issues/120
				// GDR can sometimes send phantom events with no content whatsoever. This is a bug, but its cause is
				// unclear. Ideally this is fixed in GDR and those events never get sent here, but as a workaround we can just
				// discard anything that has no Action value as garbage, to avoid it showing up in the report endpoint.
				if event.Action == "" {
					continue
				}
				s.logger.V(logging.DebugLevel).Info("updated gateway entity", "action", event.Action, "kind", event.Entity.Kind, "name", event.Entity.Name)
				eventDiff := diagnostics.NewEntityDiff(event.Diff, string(event.Action), event.Entity)
				diff.Entities = append(diff.Entities, eventDiff)
			} else {
				s.logger.Error(event.Error, "failed updating gateway entity", "action", event.Action, "kind", event.Entity.Kind, "name", event.Entity.Name)
				parsed, err := resourceErrorFromEntityAction(event)
				if err != nil {
					s.logger.Error(err, "could not parse entity update error")
				} else {
					s.resourceErrors = append(s.resourceErrors, parsed)
				}
			}
		case <-ctx.Done():
			// Release resource error lock before sending diffs to diagnostic server to prevent blocking of main procedure of updating.
			s.resourceErrorLock.Unlock()
			if diagnostic != nil && diagnostic.Diffs != nil {
				diff.Timestamp = time.Now().Format(time.RFC3339)
				diagnostic.Diffs <- diff
				s.logger.V(logging.DebugLevel).Info("recorded database update events and diff", "hash", hash)
			}
			return
		}
	}
}

func resourceErrorFromEntityAction(event diff.EntityAction) (ResourceError, error) {
	var subj any
	// GDR may produce an old only (delete), new only (create), or both (update) in an event. tags should be identical
	// but we arbitrarily pull from new.
	if event.Entity.New != nil {
		subj = event.Entity.New
	} else {
		subj = event.Entity.Old
	}
	// GDR makes frequent use of "any" for its various entity handlers. It does not use interfaces that would allow us
	// to guarantee that a particular entity does indeed have tags or similar and retrieve them. We're unlikely to
	// refactor this any time soon, so in absence of proper interface methods, we pray that the entity probably has tags,
	// which is a reasonable assumption as anything KIC can manage does. The reflect-fu here is sinister and menacing,
	// but should spit out tags unless something has gone wrong.
	reflected := reflect.Indirect(reflect.ValueOf(subj))
	if reflected.Kind() != reflect.Struct {
		// We need to fail fast here because FieldByName() will panic on non-Struct Kinds.
		return ResourceError{}, fmt.Errorf("entity %s/%s is %s, not Struct",
			event.Entity.Kind, event.Entity.Name, reflected.Kind())
	}
	tagsValue := reflected.FieldByName("Tags")
	if !tagsValue.IsValid() || tagsValue.IsZero() {
		return ResourceError{}, fmt.Errorf("entity %s/%s of type %s lacks 'Tags' field",
			event.Entity.Kind, event.Entity.Name, reflect.TypeOf(subj))
	}
	tags, ok := tagsValue.Interface().([]*string)
	if !ok {
		return ResourceError{}, fmt.Errorf("entity %s/%s Tags field is not []*string",
			event.Entity.Kind, event.Entity.Name)
	}

	actualTags := []string{}
	for _, s := range tags {
		actualTags = append(actualTags, *s)
	}

	// This omits ID, which should be available but requires similar reflect gymnastics as Tags, and probably isn't worth
	// it.
	raw := rawResourceError{
		Name: event.Entity.Name,
		Tags: actualTags,
		// /config flattened errors have a structured set of field to error reasons, whereas GDR errors are just plain
		// un-parsed admin API endpoint strings. These will often mention a field within the string, e.g.
		// schema violation (methods: cannot set 'methods' when 'protocols' is 'grpc' or 'grpcs')
		// has "methods", but we'd need to do string parsing to extract it, and we may not catch all possible error types.
		// This lazier approach just dumps the full error string as a single problem, which is probably good enough.
		Problems: map[string]string{
			fmt.Sprintf("%s:%s", event.Entity.Kind, event.Entity.Name): fmt.Sprintf("%s", event.Error),
		},
	}

	return parseRawResourceError(raw)
}

func (s *UpdateStrategyDBMode) MetricsProtocol() metrics.Protocol {
	return metrics.ProtocolDeck
}

func (s *UpdateStrategyDBMode) Type() string {
	return "DBMode"
}

func (s *UpdateStrategyDBMode) currentState(ctx context.Context) (*state.KongState, error) {
	rawState, err := dump.Get(ctx, s.client, s.dumpConfig)
	if err != nil {
		return nil, fmt.Errorf("loading configuration from kong: %w", err)
	}

	return state.Get(rawState)
}

func (s *UpdateStrategyDBMode) targetState(
	ctx context.Context,
	currentState *state.KongState,
	targetContent *file.Content,
) (*state.KongState, error) {
	rawState, err := file.Get(ctx, targetContent, file.RenderConfig{
		CurrentState: currentState,
		KongVersion:  s.version,
	}, s.dumpConfig, s.client)
	if err != nil {
		return nil, err
	}

	return state.Get(rawState)
}

// credentialMatchKey returns a string that uniquely identifies a credential within the KongState
// across syncs. It combines the consumer's ID and a sorted, canonical form of the credential's tags.
// KIC always tags credentials with GenerateTagsForObject(secret), which encodes the k8s Secret's
// namespace/name/uid — this is deterministic and survives SanitizedCopy. The consumer ID is stable
// because FillIDs covers consumers. Using these two fields (rather than the credential value such as
// Key or Username) avoids the sanitization-induced churn where SanitizedCopy randomizes key-auth Key.
func credentialMatchKey(consumerID string, tags []*string) string {
	tagStrs := make([]string, 0, len(tags))
	for _, t := range tags {
		if t != nil {
			tagStrs = append(tagStrs, *t)
		}
	}
	sort.Strings(tagStrs)
	return consumerID + "|" + strings.Join(tagStrs, ",")
}

// refillCredentialIDs reconciles credential IDs from the current state into the target state.
// When KIC builds target content, credentials carry no IDs (go-kong has no FillID for credential
// types). On the Konnect sync path, SanitizedCopy further randomizes key-auth Key values, so
// go-database-reconciler's value-based ID recovery always misses, causing delete+create churn.
// This function matches current-state credentials to target-state credentials by consumer ID + tags
// and copies the existing ID, turning churn into a stable in-place update (or no event with Fix A).
func refillCredentialIDs(currentState *state.KongState, targetState *state.KongState, logger logr.Logger) error {
	if err := refillKeyAuthIDs(currentState, targetState, logger); err != nil {
		return err
	}
	if err := refillBasicAuthIDs(currentState, targetState, logger); err != nil {
		return err
	}
	if err := refillHMACAuthIDs(currentState, targetState, logger); err != nil {
		return err
	}
	if err := refillJWTAuthIDs(currentState, targetState, logger); err != nil {
		return err
	}
	if err := refillACLGroupIDs(currentState, targetState, logger); err != nil {
		return err
	}
	if err := refillOauth2CredIDs(currentState, targetState, logger); err != nil {
		return err
	}
	return refillMTLSAuthIDs(currentState, targetState, logger)
}

type credential interface {
	state.KeyAuth |
		state.BasicAuth |
		state.HMACAuth |
		state.JWTAuth |
		state.ACLGroup |
		state.Oauth2Credential |
		state.MTLSAuth
}

// credentialOps holds type-specific operations for one credential collection, used by refillCredTypeIDs.
type credentialOps[T credential] struct {
	kind        string
	getCurrents func() ([]*T, error)
	getTargets  func() ([]*T, error)
	consumerID  func(*T) string
	id          func(*T) *string
	tags        func(*T) []*string
	setID       func(*T, *string)
	delete      func(string) error
	add         func(T) error
}

// refillCredTypeIDs is the shared implementation for all per-type refill functions.
// It builds an index of current-state IDs keyed by credentialMatchKey(consumerID, tags) and
// copies each matching ID into the target state (delete old entry + re-add with existing ID).
func refillCredTypeIDs[T credential](ops credentialOps[T], logger logr.Logger) error {
	current, err := ops.getCurrents()
	if err != nil {
		return fmt.Errorf("failed getting current %s: %w", ops.kind, err)
	}

	currentIndex := make(map[string]*string, len(current))
	for _, c := range current {
		id := ops.id(c)
		if id == nil {
			continue
		}
		k := credentialMatchKey(ops.consumerID(c), ops.tags(c))
		currentIndex[k] = id
	}

	targets, err := ops.getTargets()
	if err != nil {
		return fmt.Errorf("failed getting target %s: %w", ops.kind, err)
	}
	for _, t := range targets {
		k := credentialMatchKey(ops.consumerID(t), ops.tags(t))
		existingID, ok := currentIndex[k]
		if !ok || existingID == nil {
			continue
		}
		currentID := ops.id(t)
		if currentID != nil && *currentID == *existingID {
			continue
		}
		logger.V(logging.DebugLevel).Info("keeping ID of existing "+ops.kind, "new_id", currentID, "old_id", existingID)
		if currentID != nil {
			if err := ops.delete(*currentID); err != nil && !errors.Is(err, state.ErrNotFound) {
				return fmt.Errorf("failed deleting target %s: %w", ops.kind, err)
			}
		}
		ops.setID(t, existingID)
		if err := ops.add(*t); err != nil {
			return fmt.Errorf("failed adding target %s with refilled ID: %w", ops.kind, err)
		}
	}
	return nil
}

func refillKeyAuthIDs(currentState *state.KongState, targetState *state.KongState, logger logr.Logger) error {
	return refillCredTypeIDs(credentialOps[state.KeyAuth]{
		kind:        "key-auth",
		getCurrents: currentState.KeyAuths.GetAll,
		getTargets:  targetState.KeyAuths.GetAll,
		consumerID: func(ka *state.KeyAuth) string {
			if ka.Consumer != nil && ka.Consumer.ID != nil {
				return *ka.Consumer.ID
			}
			return ""
		},
		id:     func(ka *state.KeyAuth) *string { return ka.ID },
		tags:   func(ka *state.KeyAuth) []*string { return ka.Tags },
		setID:  func(ka *state.KeyAuth, id *string) { ka.ID = id },
		delete: targetState.KeyAuths.Delete,
		add:    targetState.KeyAuths.Add,
	}, logger)
}

func refillBasicAuthIDs(currentState *state.KongState, targetState *state.KongState, logger logr.Logger) error {
	return refillCredTypeIDs(credentialOps[state.BasicAuth]{
		kind:        "basic-auth",
		getCurrents: currentState.BasicAuths.GetAll,
		getTargets:  targetState.BasicAuths.GetAll,
		consumerID: func(ba *state.BasicAuth) string {
			if ba.Consumer != nil && ba.Consumer.ID != nil {
				return *ba.Consumer.ID
			}
			return ""
		},
		id:     func(ba *state.BasicAuth) *string { return ba.ID },
		tags:   func(ba *state.BasicAuth) []*string { return ba.Tags },
		setID:  func(ba *state.BasicAuth, id *string) { ba.ID = id },
		delete: targetState.BasicAuths.Delete,
		add:    targetState.BasicAuths.Add,
	}, logger)
}

func refillHMACAuthIDs(currentState *state.KongState, targetState *state.KongState, logger logr.Logger) error {
	return refillCredTypeIDs(credentialOps[state.HMACAuth]{
		kind:        "hmac-auth",
		getCurrents: currentState.HMACAuths.GetAll,
		getTargets:  targetState.HMACAuths.GetAll,
		consumerID: func(ha *state.HMACAuth) string {
			if ha.Consumer != nil && ha.Consumer.ID != nil {
				return *ha.Consumer.ID
			}
			return ""
		},
		id:     func(ha *state.HMACAuth) *string { return ha.ID },
		tags:   func(ha *state.HMACAuth) []*string { return ha.Tags },
		setID:  func(ha *state.HMACAuth, id *string) { ha.ID = id },
		delete: targetState.HMACAuths.Delete,
		add:    targetState.HMACAuths.Add,
	}, logger)
}

func refillJWTAuthIDs(currentState *state.KongState, targetState *state.KongState, logger logr.Logger) error {
	return refillCredTypeIDs(credentialOps[state.JWTAuth]{
		kind:        "jwt-auth",
		getCurrents: currentState.JWTAuths.GetAll,
		getTargets:  targetState.JWTAuths.GetAll,
		consumerID: func(ja *state.JWTAuth) string {
			if ja.Consumer != nil && ja.Consumer.ID != nil {
				return *ja.Consumer.ID
			}
			return ""
		},
		id:     func(ja *state.JWTAuth) *string { return ja.ID },
		tags:   func(ja *state.JWTAuth) []*string { return ja.Tags },
		setID:  func(ja *state.JWTAuth, id *string) { ja.ID = id },
		delete: targetState.JWTAuths.Delete,
		add:    targetState.JWTAuths.Add,
	}, logger)
}

func refillACLGroupIDs(currentState *state.KongState, targetState *state.KongState, logger logr.Logger) error {
	return refillCredTypeIDs(credentialOps[state.ACLGroup]{
		kind:        "acl-group",
		getCurrents: currentState.ACLGroups.GetAll,
		getTargets:  targetState.ACLGroups.GetAll,
		consumerID: func(ag *state.ACLGroup) string {
			if ag.Consumer != nil && ag.Consumer.ID != nil {
				return *ag.Consumer.ID
			}
			return ""
		},
		id:     func(ag *state.ACLGroup) *string { return ag.ID },
		tags:   func(ag *state.ACLGroup) []*string { return ag.Tags },
		setID:  func(ag *state.ACLGroup, id *string) { ag.ID = id },
		delete: targetState.ACLGroups.Delete,
		add:    targetState.ACLGroups.Add,
	}, logger)
}

func refillOauth2CredIDs(currentState *state.KongState, targetState *state.KongState, logger logr.Logger) error {
	return refillCredTypeIDs(credentialOps[state.Oauth2Credential]{
		kind:        "oauth2-cred",
		getCurrents: currentState.Oauth2Creds.GetAll,
		getTargets:  targetState.Oauth2Creds.GetAll,
		consumerID: func(oc *state.Oauth2Credential) string {
			if oc.Consumer != nil && oc.Consumer.ID != nil {
				return *oc.Consumer.ID
			}
			return ""
		},
		id:     func(oc *state.Oauth2Credential) *string { return oc.ID },
		tags:   func(oc *state.Oauth2Credential) []*string { return oc.Tags },
		setID:  func(oc *state.Oauth2Credential, id *string) { oc.ID = id },
		delete: targetState.Oauth2Creds.Delete,
		add:    targetState.Oauth2Creds.Add,
	}, logger)
}

func refillMTLSAuthIDs(currentState *state.KongState, targetState *state.KongState, logger logr.Logger) error {
	return refillCredTypeIDs(credentialOps[state.MTLSAuth]{
		kind:        "mtls-auth",
		getCurrents: currentState.MTLSAuths.GetAll,
		getTargets:  targetState.MTLSAuths.GetAll,
		consumerID: func(ma *state.MTLSAuth) string {
			if ma.Consumer != nil && ma.Consumer.ID != nil {
				return *ma.Consumer.ID
			}
			return ""
		},
		id:     func(ma *state.MTLSAuth) *string { return ma.ID },
		tags:   func(ma *state.MTLSAuth) []*string { return ma.Tags },
		setID:  func(ma *state.MTLSAuth, id *string) { ma.ID = id },
		delete: targetState.MTLSAuths.Delete,
		add:    targetState.MTLSAuths.Add,
	}, logger)
}

// refillPluginIDs keeps the plugin ID in the target state if there are already the same plugin
// (identified by plugin name and attached service, route, consumer, consumer group) in the current state.
// This prevents conflicts during the upgrade where the existing plugins have different IDs with the ID generated in building kong state.
func (s *UpdateStrategyDBMode) refillPluginIDs(currentState *state.KongState, targetState *state.KongState) error {
	plugins, err := currentState.Plugins.GetAll()
	if err != nil {
		return fmt.Errorf("failed getting plugins in current state for %s: %w", s.client.BaseRootURL(), err)
	}
	// For each existing plugin in the DB, we look for the same plugin in the target state and re-fill the ID.
	for _, existingPlugin := range plugins {
		var serviceID, routeID, consumerID, consumerGroupID string
		if existingPlugin.Service != nil && existingPlugin.Service.ID != nil {
			serviceID = *existingPlugin.Service.ID
		}
		if existingPlugin.Route != nil && existingPlugin.Route.ID != nil {
			routeID = *existingPlugin.Route.ID
		}
		if existingPlugin.Consumer != nil && existingPlugin.Consumer.ID != nil {
			consumerID = *existingPlugin.Consumer.ID
		}
		if existingPlugin.ConsumerGroup != nil && existingPlugin.ConsumerGroup.ID != nil {
			consumerGroupID = *existingPlugin.ConsumerGroup.ID
		}
		// If the same plugin is in the target state and we have filled a different ID with the existing plugin,
		// we re-fill the ID of the plugin in the target state to keep the ID the same as the existing plugin.
		targetPlugin, err := targetState.Plugins.GetByProp(*existingPlugin.Name, serviceID, routeID, consumerID, consumerGroupID)
		if err != nil {
			if !errors.Is(err, state.ErrNotFound) {
				s.logger.Error(err, "failed to get plugin with given fields in the target state")
			}
			continue
		}
		if existingPlugin.ID != nil && targetPlugin.ID != nil && *targetPlugin.ID != *existingPlugin.ID {
			s.logger.V(logging.DebugLevel).Info("Keeping ID of existing plugin",
				"plugin_name", *existingPlugin.Name, "new_plugin_id", *targetPlugin.ID, "old_plugin_id", *existingPlugin.ID,
				"service", serviceID, "route", routeID, "consumer", consumerID, "consumer_group", consumerGroupID,
			)
			// The memdb to store the target state uses `id` to identify the plugin.
			// So we need to delete the plugin with the new ID first then insert the same plugin with the same ID as the existing plugin.
			err = targetState.Plugins.Delete(*targetPlugin.ID)
			if err != nil {
				// Ignore the error if the error is ErrNotFound indicating that the plugin with the ID does not exist.
				// Otherwise, return the error and fail the update process.
				if !errors.Is(err, state.ErrNotFound) {
					s.logger.Error(err, "failed to get plugin with given ID in the target state", "id", *targetPlugin.ID)
					return err
				}
				continue
			}
			targetPlugin.ID = existingPlugin.ID
			err = targetState.Plugins.Add(*targetPlugin)
			if err != nil {
				// return error and fail the update process if we failed to add the plugin with the old ID back.
				return err
			}
		}
	}
	return nil
}
