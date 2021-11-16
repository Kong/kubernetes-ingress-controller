package admission

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	configuration "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

const (
	DefaultAdmissionWebhookCertPath = "/admission-webhook/tls.crt"
	DefaultAdmissionWebhookKeyPath  = "/admission-webhook/tls.key"
)

type ServerConfig struct {
	ListenAddr string

	CertPath string
	Cert     string

	KeyPath string
	Key     string
}

func readKeyPairFiles(certPath, keyPath string) ([]byte, []byte, error) {
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read cert from file %q: %w", certPath, err)
	}
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read key from file %q: %w", keyPath, err)
	}

	return cert, key, nil
}

func (sc *ServerConfig) toTLSConfig() (*tls.Config, error) {
	var cert, key []byte
	switch {
	case sc.CertPath == "" && sc.KeyPath == "" && sc.Cert != "" && sc.Key != "":
		cert, key = []byte(sc.Cert), []byte(sc.Key)

	case sc.CertPath != "" && sc.KeyPath != "" && sc.Cert == "" && sc.Key == "":
		var err error
		cert, key, err = readKeyPairFiles(sc.CertPath, sc.KeyPath)
		if err != nil {
			return nil, err
		}

	case sc.CertPath == "" && sc.KeyPath == "" && sc.Cert == "" && sc.Key == "":
		var err error
		cert, key, err = readKeyPairFiles(DefaultAdmissionWebhookCertPath, DefaultAdmissionWebhookKeyPath)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("either cert/key files OR cert/key values must be provided, or none")
	}

	keyPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("X509KeyPair error: %w", err)
	}
	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{keyPair},
	}, nil
}

func MakeTLSServer(config *ServerConfig, handler http.Handler) (*http.Server, error) {
	tlsConfig, err := config.toTLSConfig()
	if err != nil {
		return nil, err
	}
	return &http.Server{
		Addr:      config.ListenAddr,
		TLSConfig: tlsConfig,
		Handler:   handler,
	}, nil
}

// RequestHandler is an HTTP server that can validate Kong Ingress Controllers'
// Custom Resources using Kubernetes Admission Webhooks.
type RequestHandler struct {
	// Validator validates the entities that the k8s API-server asks
	// it the server to validate.
	Validator KongValidator

	Logger logrus.FieldLogger
}

// ServeHTTP parses AdmissionReview requests and responds back
// with the validation result of the entity.
func (a RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		a.Logger.Error("received request with empty body")
		http.Error(w, "admission review object is missing",
			http.StatusBadRequest)
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		a.Logger.Errorf("failed to read request from client: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	review := admission.AdmissionReview{}
	if err := json.Unmarshal(data, &review); err != nil {
		a.Logger.Errorf("failed to parse AdmissionReview object: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := a.handleValidation(r.Context(), *review.Request)
	if err != nil {
		a.Logger.Errorf("failed to run validation: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	review.Response = response
	data, err = json.Marshal(review)
	if err != nil {
		a.Logger.Errorf("failed to marshal response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		a.Logger.Errorf("failed to write response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var (
	consumerGVResource = meta.GroupVersionResource{
		Group:    configuration.SchemeGroupVersion.Group,
		Version:  configuration.SchemeGroupVersion.Version,
		Resource: "kongconsumers",
	}
	pluginGVResource = meta.GroupVersionResource{
		Group:    configuration.SchemeGroupVersion.Group,
		Version:  configuration.SchemeGroupVersion.Version,
		Resource: "kongplugins",
	}
	clusterPluginGVResource = meta.GroupVersionResource{
		Group:    configuration.SchemeGroupVersion.Group,
		Version:  configuration.SchemeGroupVersion.Version,
		Resource: "kongclusterplugins",
	}
	secretGVResource = meta.GroupVersionResource{
		Group:    corev1.SchemeGroupVersion.Group,
		Version:  corev1.SchemeGroupVersion.Version,
		Resource: "secrets",
	}
)

func (a RequestHandler) handleValidation(ctx context.Context, request admission.AdmissionRequest) (
	*admission.AdmissionResponse, error) {
	var response admission.AdmissionResponse

	var ok bool
	var message string
	var err error

	//nolint:exhaustive
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
			ok, message, err = a.Validator.ValidateConsumer(ctx, consumer)
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
				ok, message, err = a.Validator.ValidateConsumer(ctx, consumer)
				if err != nil {
					return nil, err
				}
			} else {
				ok = true
			}
		default:
			return nil, fmt.Errorf("unknown operation '%v'", string(request.Operation))
		}

	case pluginGVResource:
		plugin := configuration.KongPlugin{}
		deserializer := codecs.UniversalDeserializer()
		_, _, err = deserializer.Decode(request.Object.Raw,
			nil, &plugin)
		if err != nil {
			return nil, err
		}

		ok, message, err = a.Validator.ValidatePlugin(ctx, plugin)
		if err != nil {
			return nil, err
		}
	case clusterPluginGVResource:
		plugin := configuration.KongClusterPlugin{}
		deserializer := codecs.UniversalDeserializer()
		_, _, err = deserializer.Decode(request.Object.Raw,
			nil, &plugin)
		if err != nil {
			return nil, err
		}

		ok, message, err = a.Validator.ValidateClusterPlugin(ctx, plugin)
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

		// secrets are only validated on update because they must be referenced by a
		// managed consumer in order for us to validate them, and because credentials
		// validation also happens at the consumer side of the reference so a
		// credentials secret can not be referenced without being validated.
		switch request.Operation {
		case admission.Update:
			ok, message, err = a.Validator.ValidateCredential(context.Background(), secret)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown operation '%v'", string(request.Operation))
		}
	default:
		return nil, fmt.Errorf("unknown resource type to validate: %s/%s %s",
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
