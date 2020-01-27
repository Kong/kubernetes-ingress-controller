package admission

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	configuration "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	"github.com/pkg/errors"
	admission "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

// Server is an HTTP server that can validate Kong Ingress Controllers'
// Custom Resources using Kubernetes Admission Webhooks.
type Server struct {
	// Validator validates the entities that the k8s API-server asks
	// it the server to validate.
	Validator KongValidator
}

// ServeHTTP parses AdmissionReview requests and responds back
// with the validation result of the entity.
func (a Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		glog.Error("admission webhook: received request with empty body")
		http.Error(w, "admission review object is missing",
			http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error("admission webhook: reading request", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	review := admission.AdmissionReview{}
	if err := json.Unmarshal(data, &review); err != nil {
		glog.Error("admission webhook: parsing AdmissionReview", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := a.handleValidation(*review.Request)
	if err != nil {
		glog.Error("admission webhook: handling webhook", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	review.Response = response
	data, err = json.Marshal(review)
	if err != nil {
		glog.Error("admission webhook: marshaling response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		glog.Error("admission webhook: writing response", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var (
	consumerGVResource = meta.GroupVersionResource{
		Group:    configuration.SchemeGroupVersion.Group,
		Version:  configuration.SchemeGroupVersion.Version,
		Resource: "kongconsumers"}
	pluginGVResource = meta.GroupVersionResource{
		Group:    configuration.SchemeGroupVersion.Group,
		Version:  configuration.SchemeGroupVersion.Version,
		Resource: "kongplugins"}
	secretGVResource = meta.GroupVersionResource{
		Group:    corev1.SchemeGroupVersion.Group,
		Version:  corev1.SchemeGroupVersion.Version,
		Resource: "secrets"}
)

func (a Server) handleValidation(request admission.AdmissionRequest) (
	*admission.AdmissionResponse, error) {
	var response admission.AdmissionResponse

	var ok bool
	var message string
	var err error

	switch request.Resource {
	case consumerGVResource:
		consumer := configuration.KongConsumer{}
		deserializer := codecs.UniversalDeserializer()
		_, _, err = deserializer.Decode(request.Object.Raw,
			nil, &consumer)
		if err != nil {
			return nil, err
		}
		switch request.Operation {
		case admission.Create:
			ok, message, err = a.Validator.ValidateConsumer(consumer)
			if err != nil {
				return nil, err
			}
		case admission.Update:
			var oldConsumer configuration.KongConsumer
			_, _, err = deserializer.Decode(request.OldObject.Raw,
				nil, &oldConsumer)
			if err != nil {
				return nil, err
			}
			// validate only if the username is being changed
			if consumer.Username != oldConsumer.Username {
				ok, message, err = a.Validator.ValidateConsumer(consumer)
				if err != nil {
					return nil, err
				}
			} else {
				ok = true
			}
		default:
			return nil, errors.New("unknown operation '" +
				string(request.Operation) + "'")
		}

	case pluginGVResource:
		plugin := configuration.KongPlugin{}
		deserializer := codecs.UniversalDeserializer()
		_, _, err = deserializer.Decode(request.Object.Raw,
			nil, &plugin)
		if err != nil {
			return nil, err
		}

		ok, message, err = a.Validator.ValidatePlugin(plugin)
		if err != nil {
			return nil, err
		}
	case secretGVResource:
		secret := corev1.Secret{}
		deserializer := codecs.UniversalDeserializer()
		_, _, err = deserializer.Decode(request.Object.Raw,
			nil, &secret)
		if err != nil {
			return nil, err
		}
		if _, ok = secret.Data["kongCredType"]; !ok {
			// secret does not look like a credential resource in Kong
			ok = true
			break
		}

		ok, message, err = a.Validator.ValidateCredential(secret)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.Errorf("unknown resource type to validate: %s/%s %s",
			request.Resource.Group, request.Resource.Version,
			request.Resource.Resource)
	}
	if err != nil {
		return nil, err
	}
	response.UID = request.UID
	response.Allowed = ok
	response.Result = &meta.Status{
		Message: message,
	}
	if !ok {
		response.Result.Code = 400
	}
	return &response, nil
}
