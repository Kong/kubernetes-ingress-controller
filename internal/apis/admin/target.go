package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"

	adminv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/admin/v1"
)

type TargetGetter interface {
	Target() TargetInterface
}

type TargetInterface interface {
	List(url.Values, string) (*adminv1.TargetList, error)
	Get(string) (*adminv1.Target, *APIResponse)
	Create(*adminv1.Target, string) (*adminv1.Target, *APIResponse)
	Delete(string, string) error
}

type targetAPI struct {
	client APIInterface
}

func (a *targetAPI) Create(target *adminv1.Target, upstream string) (*adminv1.Target, *APIResponse) {
	rawData, err := json.Marshal(target)
	if err != nil {
		return nil, &APIResponse{err: err}
	}
	resp := a.client.RestClient().Post().
		Resource("upstreams").
		SubResource(upstream, "targets").
		Body(rawData).
		SetHeader("Content-Type", "application/json").
		Do()

	statusCode := reflect.ValueOf(resp).FieldByName("statusCode").Int()
	raw, err := resp.Raw()
	response := &APIResponse{StatusCode: int(statusCode), err: err}

	if err != nil {
		response.Raw = raw
		return nil, response
	}

	api := &adminv1.Target{}
	response.err = json.Unmarshal(raw, api)
	return api, response
}

func (a *targetAPI) Get(name string) (*adminv1.Target, *APIResponse) {
	out := &adminv1.Target{}
	err := a.client.Get(name, out)
	if err.StatusCode != http.StatusOK {
		return nil, err
	}

	return out, nil
}

func (a *targetAPI) List(params url.Values, upstream string) (*adminv1.TargetList, error) {
	if params == nil {
		params = url.Values{}
	}

	targets := &adminv1.TargetList{}
	request := a.client.
		RestClient().
		Get().
		Resource("upstreams").
		SubResource(upstream, "targets")

	for k, vals := range params {
		for _, v := range vals {
			request.Param(k, v)
		}
	}
	data, err := request.DoRaw()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, targets); err != nil {
		return nil, err
	}

	if len(targets.NextPage) > 0 {
		params.Set("offset", targets.Offset)
		result, err := a.List(params, upstream)
		if err != nil {
			return nil, err
		}
		targets.Items = append(targets.Items, result.Items...)
	}

	return targets, err
}

func (a *targetAPI) Delete(id, upstream string) error {
	return a.client.
		RestClient().
		Delete().
		RequestURI(fmt.Sprintf("/upstreams/%v/targets/%v", upstream, id)).
		Do().
		Error()
}
