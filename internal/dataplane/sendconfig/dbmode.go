package sendconfig

import (
	"context"
	"errors"
	"fmt"
	"reflect"
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
	diagnostic        *diagnostics.ClientDiagnostic
	isKonnect         bool
	logger            logr.Logger
	resourceErrors    []ResourceError
	resourceErrorLock *sync.Mutex
}

func NewUpdateStrategyDBMode(
	client *kong.Client,
	dumpConfig dump.Config,
	version semver.Version,
	concurrency int,
	diagnostic *diagnostics.ClientDiagnostic,
	logger logr.Logger,
) *UpdateStrategyDBMode {
	return &UpdateStrategyDBMode{
		client:            client,
		dumpConfig:        dumpConfig,
		version:           version,
		concurrency:       concurrency,
		diagnostic:        diagnostic,
		logger:            logger,
		resourceErrors:    []ResourceError{},
		resourceErrorLock: &sync.Mutex{},
	}
}

func NewUpdateStrategyDBModeKonnect(
	client *kong.Client,
	dumpConfig dump.Config,
	version semver.Version,
	concurrency int,
	diagnostic *diagnostics.ClientDiagnostic,
	logger logr.Logger,
) *UpdateStrategyDBMode {
	s := NewUpdateStrategyDBMode(client, dumpConfig, version, concurrency, diagnostic, logger)
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
	diagnostic *diagnostics.ClientDiagnostic,
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
			// REVIEW: Is there a better way to judge whether dumping diffs enabled than directly checking the channel?
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
