package admission

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/zapr"
	kongv1 "github.com/kong/kubernetes-configuration/api/configuration/v1"
	kongv1alpha1 "github.com/kong/kubernetes-configuration/api/configuration/v1alpha1"
	kongv1beta1 "github.com/kong/kubernetes-configuration/api/configuration/v1beta1"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
)

var decoder = codecs.UniversalDeserializer()

type KongFakeValidator struct {
	Result  bool
	Message string
	Error   error
}

func (v KongFakeValidator) ValidateConsumer(
	_ context.Context,
	_ kongv1.KongConsumer,
) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidateConsumerGroup(
	_ context.Context,
	_ kongv1beta1.KongConsumerGroup,
) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidatePlugin(
	_ context.Context,
	_ kongv1.KongPlugin,
	_ []*corev1.Secret,
) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidateClusterPlugin(
	_ context.Context,
	_ kongv1.KongClusterPlugin,
	_ []*corev1.Secret,
) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidateCredential(context.Context, corev1.Secret) (bool, string) {
	return v.Result, v.Message
}

func (v KongFakeValidator) ValidateGateway(_ context.Context, _ gatewayapi.Gateway) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidateHTTPRoute(_ context.Context, _ gatewayapi.HTTPRoute) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidateIngress(_ context.Context, _ netv1.Ingress) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidateVault(_ context.Context, _ kongv1alpha1.KongVault) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidateCustomEntity(_ context.Context, _ kongv1alpha1.KongCustomEntity) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func (v KongFakeValidator) ValidateService(_ context.Context, _ corev1.Service) (bool, string, error) {
	return v.Result, v.Message, v.Error
}

func TestServeHTTPBasic(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	server := RequestHandler{
		Validator: KongFakeValidator{},
		Logger:    zapr.NewLogger(zap.NewNop()),
	}
	handler := http.HandlerFunc(server.ServeHTTP)

	req, err := http.NewRequest("POST", "", nil)
	assert.Nil(err)
	handler.ServeHTTP(res, req)
	assert.Equal(400, res.Code)
	assert.Equal("Admission review object is missing\n",
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
			wantSuccessResponse admissionv1.AdmissionResponse
			wantFailureMessage  string
		}{
			{
				name:               "request with present empty body",
				wantRespCode:       http.StatusBadRequest,
				wantFailureMessage: "EOF\n",
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
				wantSuccessResponse: admissionv1.AdmissionResponse{
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
				wantSuccessResponse: admissionv1.AdmissionResponse{
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
				wantSuccessResponse: admissionv1.AdmissionResponse{
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
				wantSuccessResponse: admissionv1.AdmissionResponse{
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
				wantSuccessResponse: admissionv1.AdmissionResponse{
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
					Logger:    zapr.NewLogger(zap.NewNop()),
				}
				handler := http.HandlerFunc(server.ServeHTTP)

				// act
				req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tt.reqBody)))
				assert.Nil(err)
				handler.ServeHTTP(res, req)

				// assert
				assert.Equal(tt.wantRespCode, res.Code)
				if tt.wantRespCode == http.StatusOK {
					var review admissionv1.AdmissionReview
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
