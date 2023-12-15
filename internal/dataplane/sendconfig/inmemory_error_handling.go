package sendconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kong/kubernetes-ingress-controller/v2/internal/util"
)

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
type ConfigErrorFields struct{}

// FlatEntityError represents a single Kong entity with one or more invalid fields.
type FlatEntityError struct {
	Type   string           `json:"entity_type,omitempty" yaml:"entity_type,omitempty"`
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

// parseFlatEntityErrors takes a Kong /config error response body and parses its "fields.flattened_errors" value
// into errors associated with Kubernetes resources.
func parseFlatEntityErrors(body []byte) ([]FlatEntityError, error) {
	// Directly return here to avoid the misleading "could not unmarshal config" message appear in logs.
	if len(body) == 0 {
		return nil, nil
	}

	var configError ConfigError

	err := json.Unmarshal(body, &configError)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config error: %w", err)
	}

	return configError.Flattened, nil
}

func ResourceErrorsFromEntityErrors(entityErrors []FlatEntityError, log logrus.FieldLogger) []ResourceError {
	var resourceErrors []ResourceError
	for _, ee := range entityErrors {
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
					"field": p.Field,
				}).Error("entity has both single and array errors for field")

				continue
			}
			if len(p.Message) > 0 {
				raw.Problems[p.Field] = p.Message
			}
			if len(p.Messages) > 0 {
				for i, message := range p.Messages {
					if len(message) > 0 {
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

	return resourceErrors
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
