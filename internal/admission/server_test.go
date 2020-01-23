package admission

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	configuration "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	admission "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

var decoder = codecs.UniversalDeserializer()

type KongFakeValidator struct {
	Result  bool
	Message string
	Error   error
}

func (v KongFakeValidator) ValidateConsumer(
	consumer configuration.KongConsumer) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidatePlugin(
	k8sPlugin configuration.KongPlugin) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidateCredential(
	secret corev1.Secret) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func TestServeHTTPBasic(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	server := Server{
		Validator: KongFakeValidator{},
	}
	handler := http.HandlerFunc(server.ServeHTTP)

	req, err := http.NewRequest("POST", "", nil)
	assert.Nil(err)
	handler.ServeHTTP(res, req)
	assert.Equal(400, res.Code)
	assert.Equal("admission review object is missing\n",
		res.Body.String())
}

func TestValidateKongConsumer(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	server := Server{
		Validator: KongFakeValidator{
			Result: true,
		},
	}
	handler := http.HandlerFunc(server.ServeHTTP)
	// TODO how to marshal k8s object to correct JSON?
	body := `
{
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "b2df61dd-ab5b-4cb4-9be0-878533c83892",
    "resource": {
      "group": "configuration.konghq.com",
      "version": "v1",
      "resource": "kongconsumers"
    },
    "object": {
      "apiVersion": "configuration.konghq.com/v1",
      "kind": "KongConsumer"
    },
	"operation": "CREATE"
  }
}
	`
	req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(body)))
	assert.Nil(err)
	handler.ServeHTTP(res, req)
	assert.Equal(200, res.Code)
	var review admission.AdmissionReview
	_, _, err = decoder.Decode([]byte(res.Body.String()), nil, &review)
	assert.Nil(err)
	assert.Equal("b2df61dd-ab5b-4cb4-9be0-878533c83892",
		string(review.Response.UID))
	assert.True(review.Response.Allowed)
}

func TestValidateKongConsumerOnUsernameChange(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	server := Server{
		Validator: KongFakeValidator{
			Result: true,
		},
	}
	handler := http.HandlerFunc(server.ServeHTTP)
	// TODO how to marshal k8s object to correct JSON?
	body := `
{
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "b2df61dd-ab5b-4cb4-9be0-878533c83892",
    "resource": {
      "group": "configuration.konghq.com",
      "version": "v1",
      "resource": "kongconsumers"
    },
    "object": {
      "apiVersion": "configuration.konghq.com/v1",
      "kind": "KongConsumer",
	  "username":"foo"
    },
    "oldObject": {
      "apiVersion": "configuration.konghq.com/v1",
      "kind": "KongConsumer",
	  "username":"bar"
    },
	"operation": "UPDATE"
  }
}
	`
	req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(body)))
	assert.Nil(err)
	handler.ServeHTTP(res, req)
	assert.Equal(200, res.Code)
	var review admission.AdmissionReview
	_, _, err = decoder.Decode([]byte(res.Body.String()), nil, &review)
	assert.Nil(err)
	assert.Equal("b2df61dd-ab5b-4cb4-9be0-878533c83892",
		string(review.Response.UID))
	assert.True(review.Response.Allowed)
}

func TestValidateKongConsumerOnEqualUpdate(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	server := Server{
		Validator: KongFakeValidator{
			Result: true,
		},
	}
	handler := http.HandlerFunc(server.ServeHTTP)
	// TODO how to marshal k8s object to correct JSON?
	body := `
{
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "b2df61dd-ab5b-4cb4-9be0-878533c83892",
    "resource": {
      "group": "configuration.konghq.com",
      "version": "v1",
      "resource": "kongconsumers"
    },
    "object": {
      "apiVersion": "configuration.konghq.com/v1",
      "kind": "KongConsumer",
	  "username":"foo"
    },
    "oldObject": {
      "apiVersion": "configuration.konghq.com/v1",
      "kind": "KongConsumer",
	  "username":"foo"
    },
	"operation": "UPDATE"
  }
}
	`
	req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(body)))
	assert.Nil(err)
	handler.ServeHTTP(res, req)
	assert.Equal(200, res.Code)
	var review admission.AdmissionReview
	_, _, err = decoder.Decode([]byte(res.Body.String()), nil, &review)
	assert.Nil(err)
	assert.Equal("b2df61dd-ab5b-4cb4-9be0-878533c83892",
		string(review.Response.UID))
	assert.True(review.Response.Allowed)
}

func TestValidateKongConsumerInvalid(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	server := Server{
		Validator: KongFakeValidator{
			Result:  false,
			Message: "consumer is not valid",
		},
	}
	handler := http.HandlerFunc(server.ServeHTTP)
	body := `
{
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "b2df61dd-ab5b-4cb4-9be0-878533c83892",
    "resource": {
      "group": "configuration.konghq.com",
      "version": "v1",
      "resource": "kongconsumers"
    },
    "object": {
      "apiVersion": "configuration.konghq.com/v1",
      "kind": "KongConsumer"
    },
	"operation": "CREATE"
  }
}
	`
	req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(body)))
	assert.Nil(err)
	handler.ServeHTTP(res, req)
	assert.Equal(200, res.Code)
	var review admission.AdmissionReview
	_, _, err = decoder.Decode([]byte(res.Body.String()), nil, &review)
	assert.Nil(err)
	assert.Equal("b2df61dd-ab5b-4cb4-9be0-878533c83892",
		string(review.Response.UID))
	assert.False(review.Response.Allowed)
	assert.Equal("consumer is not valid", review.Response.Result.Message)
	assert.Equal(int32(400), review.Response.Result.Code)
}

func TestValidateKongConsumerOnError(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	server := Server{
		Validator: KongFakeValidator{
			Error: errors.New("error making API call to kong"),
		},
	}
	handler := http.HandlerFunc(server.ServeHTTP)
	body := `
{
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "b2df61dd-ab5b-4cb4-9be0-878533c83892",
    "resource": {
      "group": "configuration.konghq.com",
      "version": "v1",
      "resource": "kongconsumers"
    },
    "object": {
      "apiVersion": "configuration.konghq.com/v1",
      "kind": "KongConsumer"
    },
	"operation": "CREATE"
  }
}
	`
	req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(body)))
	assert.Nil(err)
	handler.ServeHTTP(res, req)
	assert.Equal(500, res.Code)
	assert.Equal("error making API call to kong\n", res.Body.String())
}

func TestValidateKongConsumerOnUsernameChangeError(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	server := Server{
		Validator: KongFakeValidator{
			Error: errors.New("error making API call to kong"),
		},
	}
	handler := http.HandlerFunc(server.ServeHTTP)
	body := `
{
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "b2df61dd-ab5b-4cb4-9be0-878533c83892",
    "resource": {
      "group": "configuration.konghq.com",
      "version": "v1",
      "resource": "kongconsumers"
    },
    "object": {
      "apiVersion": "configuration.konghq.com/v1",
      "kind": "KongConsumer",
	  "username":"foo"
    },
    "oldObject": {
      "apiVersion": "configuration.konghq.com/v1",
      "kind": "KongConsumer",
	  "username":"bar"
    },
	"operation": "UPDATE"
  }
}
	`
	req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(body)))
	assert.Nil(err)
	handler.ServeHTTP(res, req)
	assert.Equal(500, res.Code)
	assert.Equal("error making API call to kong\n", res.Body.String())
}

func TestUnknownResource(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	server := Server{
		Validator: KongFakeValidator{
			Result:  false,
			Message: "consumer is not valid",
		},
	}
	handler := http.HandlerFunc(server.ServeHTTP)
	body := `
{
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "b2df61dd-ab5b-4cb4-9be0-878533c83892",
    "resource": {
      "group": "configuration.konghq.com",
      "version": "v1",
      "resource": "kongunknown"
    },
    "object": {
      "apiVersion": "configuration.konghq.com/v1",
      "kind": "KongConsumer"
    },
	"operation": "CREATE"
  }
}
	`
	req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(body)))
	assert.Nil(err)
	handler.ServeHTTP(res, req)
	assert.Equal(500, res.Code)
	assert.Equal("unknown resource type to validate: configuration.konghq.com/v1 kongunknown\n", res.Body.String())
}

func TestValidateKongPlugin(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	server := Server{
		Validator: KongFakeValidator{
			Result: true,
		},
	}
	handler := http.HandlerFunc(server.ServeHTTP)
	body := `
{
  "kind": "AdmissionReview",
  "apiVersion": "admission.k8s.io/v1beta1",
  "request": {
    "uid": "b2df61dd-ab5b-4cb4-9be0-878533c83892",
    "resource": {
      "group": "configuration.konghq.com",
      "version": "v1",
      "resource": "kongplugins"
    },
    "object": {
      "apiVersion": "configuration.konghq.com/v1",
      "kind": "KongPlugin"
    }
  }
}
	`
	req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(body)))
	assert.Nil(err)
	handler.ServeHTTP(res, req)
	assert.Equal(200, res.Code)
	var review admission.AdmissionReview
	_, _, err = decoder.Decode([]byte(res.Body.String()), nil, &review)
	assert.Nil(err)
	assert.Equal("b2df61dd-ab5b-4cb4-9be0-878533c83892",
		string(review.Response.UID))
	assert.True(review.Response.Allowed)
}
