package admission

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configuration "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

var decoder = codecs.UniversalDeserializer()

type KongFakeValidator struct {
	Result  bool
	Message string
	Error   error
}

func (v KongFakeValidator) ValidateConsumer(_ context.Context,
	consumer configuration.KongConsumer) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidatePlugin(_ context.Context,
	k8sPlugin configuration.KongPlugin) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidateClusterPlugin(_ context.Context,
	k8sPlugin configuration.KongClusterPlugin) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidateCredential(ctx context.Context, secret corev1.Secret) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func TestServeHTTPBasic(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	server := RequestHandler{
		Validator: KongFakeValidator{},
		Logger:    logrus.New(),
	}
	handler := http.HandlerFunc(server.ServeHTTP)

	req, err := http.NewRequest("POST", "", nil)
	assert.Nil(err)
	handler.ServeHTTP(res, req)
	assert.Equal(400, res.Code)
	assert.Equal("admission review object is missing\n",
		res.Body.String())
}

func TestValidationWebhook(t *testing.T) {
	for _, apiVersion := range []string{
		"admission.k8s.io/v1beta1",
		"admission.k8s.io/v1",
	} {
		for _, tt := range []struct {
			name      string
			reqBody   string
			validator KongValidator

			wantRespCode        int
			wantSuccessResponse admission.AdmissionResponse
			wantFailureMessage  string
		}{
			{
				name:               "request with present empty body",
				wantRespCode:       http.StatusBadRequest,
				wantFailureMessage: "unexpected end of JSON input\n",
			},
			{
				name: "validate kong consumer",
				reqBody: dedent.Dedent(`
					{
						"kind": "AdmissionReview",
						"apiVersion": "` + apiVersion + `",
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
					}`),
				validator:    KongFakeValidator{Result: true},
				wantRespCode: http.StatusOK,
				wantSuccessResponse: admission.AdmissionResponse{
					UID:     "b2df61dd-ab5b-4cb4-9be0-878533c83892",
					Allowed: true,
					Result:  &metav1.Status{},
				},
			},
			{
				name: "validate kong consumer on username change",
				reqBody: dedent.Dedent(`
					{
						"kind": "AdmissionReview",
						"apiVersion": "` + apiVersion + `",
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
				`),
				validator:    KongFakeValidator{Result: true},
				wantRespCode: http.StatusOK,
				wantSuccessResponse: admission.AdmissionResponse{
					UID:     "b2df61dd-ab5b-4cb4-9be0-878533c83892",
					Allowed: true,
					Result:  &metav1.Status{},
				},
			},
			{
				name: "validate kong consumer on equal update",
				reqBody: dedent.Dedent(`
					{
						"kind": "AdmissionReview",
						"apiVersion": "` + apiVersion + `",
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
					}`),
				validator:    KongFakeValidator{Result: true},
				wantRespCode: http.StatusOK,
				wantSuccessResponse: admission.AdmissionResponse{
					UID:     "b2df61dd-ab5b-4cb4-9be0-878533c83892",
					Allowed: true,
					Result:  &metav1.Status{},
				},
			},
			{
				name: "validate kong consumer invalid",
				reqBody: dedent.Dedent(`
					{
						"kind": "AdmissionReview",
						"apiVersion": "` + apiVersion + `",
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
					}`),
				validator:    KongFakeValidator{Result: false, Message: "consumer is not valid"},
				wantRespCode: http.StatusOK,
				wantSuccessResponse: admission.AdmissionResponse{
					UID:     "b2df61dd-ab5b-4cb4-9be0-878533c83892",
					Allowed: false,
					Result: &metav1.Status{
						Code:    http.StatusBadRequest,
						Message: "consumer is not valid",
					},
				},
			},
			{
				name: "kong consumer validator error",
				reqBody: dedent.Dedent(`
					{
						"kind": "AdmissionReview",
						"apiVersion": "` + apiVersion + `",
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
					}`),
				validator:          KongFakeValidator{Error: errors.New("error making API call to kong")},
				wantRespCode:       http.StatusInternalServerError,
				wantFailureMessage: "error making API call to kong\n",
			},
			{
				name: "kong consumer validator error on username change",
				reqBody: dedent.Dedent(`
					{
						"kind": "AdmissionReview",
						"apiVersion": "` + apiVersion + `",
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
					}`),
				validator:          KongFakeValidator{Error: errors.New("error making API call to kong")},
				wantRespCode:       http.StatusInternalServerError,
				wantFailureMessage: "error making API call to kong\n",
			},
			{
				name: "unknown resource",
				reqBody: dedent.Dedent(`
					{
						"kind": "AdmissionReview",
						"apiVersion": "` + apiVersion + `",
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
					}`),
				validator:          KongFakeValidator{Result: false, Message: "consumer is not valid"},
				wantRespCode:       http.StatusInternalServerError,
				wantFailureMessage: "unknown resource type to validate: configuration.konghq.com/v1 kongunknown\n",
			},
			{
				name: "validate kong plugin",
				reqBody: dedent.Dedent(`
					{
						"kind": "AdmissionReview",
						"apiVersion": "` + apiVersion + `",
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
					}`),
				validator:    KongFakeValidator{Result: true},
				wantRespCode: http.StatusOK,
				wantSuccessResponse: admission.AdmissionResponse{
					UID:     "b2df61dd-ab5b-4cb4-9be0-878533c83892",
					Allowed: true,
					Result:  &metav1.Status{},
				},
			},
		} {
			t.Run(fmt.Sprintf("%s/%s", apiVersion, tt.name), func(t *testing.T) {
				// arrange
				assert := assert.New(t)
				res := httptest.NewRecorder()
				server := RequestHandler{
					Validator: tt.validator,
					Logger:    logrus.New(),
				}
				handler := http.HandlerFunc(server.ServeHTTP)

				// act
				req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tt.reqBody)))
				assert.Nil(err)
				handler.ServeHTTP(res, req)

				// assert
				assert.Equal(tt.wantRespCode, res.Code)
				if tt.wantRespCode == http.StatusOK {
					var review admission.AdmissionReview
					_, _, err = decoder.Decode(res.Body.Bytes(), nil, &review)
					assert.Nil(err)
					assert.EqualValues(&tt.wantSuccessResponse, review.Response)
				} else {
					assert.Equal(res.Body.String(), tt.wantFailureMessage)
				}
			})
		}
	}
}
