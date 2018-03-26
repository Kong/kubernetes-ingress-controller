package admin

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"

	adminv1 "github.com/kong/ingress-controller/internal/apis/admin/v1"
	"k8s.io/client-go/rest"
)

type CredentialGetter interface {
	Credentials() CredentialInterface
}

type CredentialInterface interface {
	List(string, url.Values) (*adminv1.CredentialList, error)
	GetByType(string, string) (*adminv1.Credential, *APIResponse)
	CreateByType(map[string]interface{}, string, string) *APIResponse
}

type credentialAPI struct {
	client rest.Interface
}

// curl -X POST http://kong:8001/consumers/{consumer}/{name} -d ''
func (a *credentialAPI) CreateByType(obj map[string]interface{}, consumer, name string) *APIResponse {
	rawData, err := json.Marshal(obj)
	if err != nil {
		return &APIResponse{err: err}
	}
	resp := a.client.Post().
		Resource("consumers").
		SubResource(consumer, name).
		Body(rawData).
		SetHeader("Content-Type", "application/json").
		Do()

	statusCode := reflect.ValueOf(resp).FieldByName("statusCode").Int()
	raw, err := resp.Raw()
	response := &APIResponse{StatusCode: int(statusCode), err: err}
	if err != nil {
		response.Raw = raw
	}

	return response
}

func (a *credentialAPI) GetByType(id, name string) (*adminv1.Credential, *APIResponse) {
	resp := a.client.Get().
		Resource(name).
		Name(id).
		Do()
	statusCode := reflect.ValueOf(resp).FieldByName("statusCode").Int()
	raw, err := resp.Raw()
	response := &APIResponse{StatusCode: int(statusCode), err: err}
	if err != nil {
		response.Raw = raw
		return nil, response
	}

	out := &adminv1.Credential{}
	response.err = json.Unmarshal(raw, out)
	return out, response

}

func (a *credentialAPI) List(name string, params url.Values) (*adminv1.CredentialList, error) {
	plural := fmt.Sprintf("%vs", name)
	credentials := &adminv1.CredentialList{}
	request := a.client.Get().Resource(plural)
	for k, vals := range params {
		for _, v := range vals {
			request.Param(k, v)
		}
	}

	data, err := request.DoRaw()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, credentials); err != nil {
		return nil, err
	}

	if len(credentials.NextPage) > 0 {
		params.Add("offset", credentials.Offset)
		result, err := a.List(name, params)
		if err != nil {
			return nil, err
		}
		credentials.Items = append(credentials.Items, result.Items...)
	}

	return credentials, err
}
