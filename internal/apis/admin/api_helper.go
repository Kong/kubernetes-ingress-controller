package admin

import (
	"encoding/json"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

type APIInterface interface {
	Get(string, interface{}) *APIResponse
	Create(interface{}, interface{}) *APIResponse
	Patch(string, interface{}, interface{}) *APIResponse
	Delete(string) error

	RestClient() rest.Interface
}

type apiClient struct {
	client   rest.Interface
	resource *metav1.APIResource
}

func (a *apiClient) RestClient() rest.Interface {
	return a.client
}

func (a *apiClient) Patch(id string, obj interface{}, out interface{}) *APIResponse {
	rawData, err := json.Marshal(obj)
	if err != nil {
		return &APIResponse{err: err}
	}
	resp := a.client.Patch(types.JSONPatchType).
		Resource(a.resource.Name).
		SubResource(id).
		Body(rawData).
		SetHeader("Content-Type", "application/json").
		Do()

	statusCode := reflect.ValueOf(resp).FieldByName("statusCode").Int()
	raw, err := resp.Raw()
	response := &APIResponse{StatusCode: int(statusCode), err: err}

	if err != nil {
		response.Raw = raw
		return response
	}

	response.err = json.Unmarshal(raw, out)
	return response
}

func (a *apiClient) Create(obj, out interface{}) *APIResponse {
	rawData, err := json.Marshal(obj)
	if err != nil {
		return &APIResponse{err: err}
	}
	resp := a.client.Post().
		Resource(a.resource.Name).
		Body(rawData).
		SetHeader("Content-Type", "application/json").
		Do()

	statusCode := reflect.ValueOf(resp).FieldByName("statusCode").Int()
	raw, err := resp.Raw()
	response := &APIResponse{StatusCode: int(statusCode), err: err}
	if err != nil {
		response.Raw = raw
		return response
	}

	response.err = json.Unmarshal(raw, out)
	return response
}

func (a *apiClient) Get(name string, out interface{}) *APIResponse {
	resp := a.client.Get().
		Resource(a.resource.Name).
		Name(name).
		Do()
	statusCode := reflect.ValueOf(resp).FieldByName("statusCode").Int()
	raw, err := resp.Raw()
	response := &APIResponse{StatusCode: int(statusCode), err: err}
	if err != nil {
		response.Raw = raw
		return response
	}

	response.err = json.Unmarshal(raw, out)
	return response
}

func (a *apiClient) Delete(id string) error {
	return a.client.Delete().
		Resource(a.resource.Name).
		Name(id).
		Do().
		Error()
}
