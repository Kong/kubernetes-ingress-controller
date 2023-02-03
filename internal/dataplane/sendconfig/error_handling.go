package sendconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"

	deckutils "github.com/kong/deck/utils"
	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/metrics"
)

// ResourceError is a Kong configuration error associated with a Kubernetes resource.
type ResourceError struct {
	Name       string
	Namespace  string
	Kind       string
	APIVersion string
	UID        string
	Problems   map[string]string
}

// rawResourceError is a Kong configuration error associated with a Kubernetes resource with Kubernetes metadata stored
// in raw Kong entity tags.
type rawResourceError struct {
	Name     string
	ID       string
	Tags     []string
	Problems map[string]string
}

// ConfigError is an error response from Kong's DB-less /config endpoint.
type ConfigError struct {
	Code      int               `json:"code,omitempty" yaml:"code,omitempty"`
	Flattened []FlatEntityError `json:"flattened_errors,omitempty" yaml:"flattened_errors,omitempty"`
	Message   string            `json:"message,omitempty" yaml:"message,omitempty"`
	Name      string            `json:"name,omitempty" yaml:"name,omitempty"`
}

// ConfigErrorFields is the structure under the "fields" key in a /config error response.
type ConfigErrorFields struct {
}

// FlatEntityError represents a single Kong entity with one or more invalid fields.
type FlatEntityError struct {
	Name   string           `json:"entity_name,omitempty" yaml:"entity_name,omitempty"`
	ID     string           `json:"entity_id,omitempty" yaml:"entity_id,omitempty"`
	Tags   []string         `json:"entity_tags,omitempty" yaml:"entity_tags,omitempty"`
	Errors []FlatFieldError `json:"errors,omitempty" yaml:"errors,omitempty"`
}

// FlatFieldError represents an error for a single field within a Kong entity.
type FlatFieldError struct {
	Field string `json:"field,omitempty" yaml:"field,omitempty"`
	// Message is the error associated with Field for single-value fields.
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	// Messages are the errors associated with Field for multi-value fields. The array index in Messages matches the
	// array index in the input.
	Messages []string `json:"messages,omitempty" yaml:"messages,omitempty"`
}

// parseFlatEntityErrors takes a Kong /config error response body and parses its "fields.flattened" value into errors
// associated with Kubernetes resources.
func parseFlatEntityErrors(body []byte, log logrus.FieldLogger) ([]ResourceError, error) {
	var resourceErrors []ResourceError
	var configError ConfigError
	err := json.Unmarshal(body, &configError)
	if err != nil {
		return resourceErrors, fmt.Errorf("could not unmarshal config error: %w", err)
	}
	for _, ee := range configError.Flattened {
		raw := rawResourceError{
			Name:     ee.Name,
			ID:       ee.ID,
			Tags:     ee.Tags,
			Problems: map[string]string{},
		}
		for _, p := range ee.Errors {
			if len(p.Message) > 0 && len(p.Messages) > 0 {
				log.WithFields(logrus.Fields{
					"name":  ee.Name,
					"field": p.Field}).Error("entity has both single and array errors for field")
				continue
			}
			if len(p.Message) > 0 {
				raw.Problems[p.Field] = p.Message
			}
			if len(p.Messages) > 0 {
				for i, message := range p.Messages {
					if len(message) > 0 { // TODO how are the nulls treated?
						raw.Problems[fmt.Sprintf("%s[%d]", p.Field, i)] = message
					}
				}
			}
		}
		parsed, err := parseRawResourceError(raw)
		if err != nil {
			log.WithError(err).WithField("name", ee.Name).Error("entity tags missing fields")
			continue
		}
		resourceErrors = append(resourceErrors, parsed)
	}
	return resourceErrors, nil
}

// parseRawResourceError takes a raw resource error and parses its tags into Kubernetes metadata. If critical tags are
// missing, it returns an error indicating the missing tag.
func parseRawResourceError(raw rawResourceError) (ResourceError, error) {
	re := ResourceError{}
	re.Problems = raw.Problems
	var gvk schema.GroupVersionKind
	for _, tag := range raw.Tags {
		if strings.HasPrefix(tag, util.K8sNameTagPrefix) {
			re.Name = strings.TrimPrefix(tag, util.K8sNameTagPrefix)
		}
		if strings.HasPrefix(tag, util.K8sNamespaceTagPrefix) {
			re.Namespace = strings.TrimPrefix(tag, util.K8sNamespaceTagPrefix)
		}
		if strings.HasPrefix(tag, util.K8sKindTagPrefix) {
			gvk.Kind = strings.TrimPrefix(tag, util.K8sKindTagPrefix)
		}
		if strings.HasPrefix(tag, util.K8sVersionTagPrefix) {
			gvk.Version = strings.TrimPrefix(tag, util.K8sVersionTagPrefix)
		}
		// this will not set anything for core resources
		if strings.HasPrefix(tag, util.K8sGroupTagPrefix) {
			gvk.Group = strings.TrimPrefix(tag, util.K8sGroupTagPrefix)
		}
		if strings.HasPrefix(tag, util.K8sUIDTagPrefix) {
			re.UID = strings.TrimPrefix(tag, util.K8sUIDTagPrefix)
		}
	}
	re.APIVersion, re.Kind = gvk.ToAPIVersionAndKind()
	if re.Name == "" {
		return re, fmt.Errorf("no name")
	}
	if re.Namespace == "" {
		return re, fmt.Errorf("no namespace")
	}
	if re.Kind == "" {
		return re, fmt.Errorf("no kind")
	}
	if re.UID == "" {
		return re, fmt.Errorf("no uid")
	}
	return re, nil
}

// deckConfigConflictError is an error used to wrap deck config conflict errors returned from deck functions
// transforming KongRawState to KongState (e.g. state.Get, dump.Get).
type deckConfigConflictError struct {
	err error
}

func (e deckConfigConflictError) Error() string {
	return e.err.Error()
}

func (e deckConfigConflictError) Is(target error) bool {
	_, ok := target.(deckConfigConflictError)
	return ok
}

func (e deckConfigConflictError) Unwrap() error {
	return e.err
}

// pushFailureReason extracts config push failure reason from an error returned from onUpdateInMemoryMode or onUpdateDBMode.
func pushFailureReason(err error) string {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return metrics.FailureReasonNetwork
	}

	if isConflictErr(err) {
		return metrics.FailureReasonConflict
	}

	return metrics.FailureReasonOther
}

func isConflictErr(err error) bool {
	var apiErr *kong.APIError
	if errors.As(err, &apiErr) && apiErr.Code() == http.StatusConflict ||
		errors.Is(err, deckConfigConflictError{}) {
		return true
	}

	var deckErrArray deckutils.ErrArray
	if errors.As(err, &deckErrArray) {
		for _, err := range deckErrArray.Errors {
			if isConflictErr(err) {
				return true
			}
		}
	}

	return false
}
