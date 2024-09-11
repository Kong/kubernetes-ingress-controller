package sendconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/util"
	kongv1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1alpha1"
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
	// Name is the name of the Kong entity.
	Name string `json:"entity_name,omitempty" yaml:"entity_name,omitempty"`

	// ID is the ID of the Kong entity.
	ID string `json:"entity_id,omitempty" yaml:"entity_id,omitempty"`

	// Tags are the tags of the Kong entity.
	Tags []string `json:"entity_tags,omitempty" yaml:"entity_tags,omitempty"`

	// Type is the type of the Kong entity.
	Type string `json:"entity_type,omitempty" yaml:"entity_type,omitempty"`

	// Errors are the errors associated with the Kong entity.
	Errors []FlatError `json:"errors,omitempty" yaml:"errors,omitempty"`
}

// FlatErrorType tells whether a FlatError is associated with a single field or a whole entity.
type FlatErrorType string

const (
	// FlatErrorTypeField is an error associated with a single field.
	FlatErrorTypeField FlatErrorType = "field"

	// FlatErrorTypeEntity is an error associated with a whole entity.
	FlatErrorTypeEntity FlatErrorType = "entity"
)

// FlatError represents an error for a single field within a Kong entity or a whole entity.
type FlatError struct {
	// Field is the name of the entity's field that has an error.
	// Optional: Field can be empty if the error is associated with the whole entity.
	Field string `json:"field,omitempty" yaml:"field,omitempty"`

	// Message is the error associated with Field (for single-value fields) or with a whole entity when Type is "entity".
	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	// Messages are the errors associated with Field for multi-value fields. The array index in Messages matches the
	// array index in the input.
	Messages []string `json:"messages,omitempty" yaml:"messages,omitempty"`

	// Type tells whether the error is associated with a single field or a whole entity.
	Type FlatErrorType `json:"type,omitempty" yaml:"type,omitempty"`
}

// parseFlatEntityErrors takes a Kong /config error response body and parses its "fields.flattened_errors" value
// into errors associated with Kubernetes resources.
func parseFlatEntityErrors(body []byte, logger logr.Logger) ([]ResourceError, error) {
	// Directly return here to avoid the misleading "could not unmarshal config" message appear in logs.
	if len(body) == 0 {
		return nil, nil
	}

	var resourceErrors []ResourceError //nolint:prealloc
	var configError ConfigError

	err := json.Unmarshal(body, &configError)
	if err != nil {
		// we _should_ arguably be able to parse the "message" field into a ConfigError even if we can't parse a full set
		// of flattened errors, but for some reason those incomplete errors still don't unmarshal. as a fallback, try and
		// yank the field out via a more basic unmarshal target
		fallback := map[string]interface{}{}
		if fallbackErr := json.Unmarshal(body, &fallback); fallbackErr == nil {
			if message, ok := fallback["message"]; ok {
				logger.Error(nil, "Could not fully parse config error", "message", message)
			}
		}
		return resourceErrors, NewResponseParsingError(body)
	}
	if len(configError.Flattened) == 0 {
		if len(configError.Message) > 0 {
			logger.Error(nil, "Config error missing per-resource errors", "message", configError.Message)
		} else {
			logger.Error(nil, "Config error missing per-resource and message", "message", configError.Message)
		}
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
				logger.Error(nil, "Entity has both single and array errors for field",
					"name", ee.Name, "field", p.Field)
				continue
			}
			if len(p.Message) > 0 {
				switch p.Type {
				case FlatErrorTypeField:
					// If the error is associated with a single field, store it in the map under the field name.
					raw.Problems[p.Field] = p.Message
				case FlatErrorTypeEntity:
					// If the error is associated with a whole entity, store it in the map under the entity type and name.
					raw.Problems[fmt.Sprintf("%s:%s", ee.Type, ee.Name)] = p.Message
				}
			}
			for i, message := range p.Messages {
				if len(message) > 0 {
					raw.Problems[fmt.Sprintf("%s[%d]", p.Field, i)] = message
				}
			}
		}
		parsed, err := parseRawResourceError(raw)
		if err != nil {
			logger.Error(err, "Entity tags missing fields", "name", ee.Name)
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
	if re.Namespace == "" && !gvkIsClusterScoped(gvk) {
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

func gvkIsClusterScoped(gvk schema.GroupVersionKind) bool {
	if gvk.Group == kongv1.GroupVersion.Group && gvk.Version == kongv1.GroupVersion.Version {
		return gvk.Kind == "KongClusterPlugin" || gvk.Kind == "KongLicense"
	}
	if gvk.Group == kongv1alpha1.GroupVersion.Group && gvk.Version == kongv1alpha1.GroupVersion.Version {
		return gvk.Kind == "KongVault"
	}
	return false
}
