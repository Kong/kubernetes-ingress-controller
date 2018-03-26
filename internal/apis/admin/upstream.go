package admin

import (
	"encoding/json"
	"net/url"

	adminv1 "github.com/kong/ingress-controller/internal/apis/admin/v1"
)

type UpstreamGetter interface {
	Upstreams() UpstreamInterface
}

type UpstreamInterface interface {
	List(params url.Values) (*adminv1.UpstreamList, error)
	Get(name string) (*adminv1.Upstream, *APIResponse)
	Create(route *adminv1.Upstream) (*adminv1.Upstream, *APIResponse)
	Delete(name string) error
}

type upstreamAPI struct {
	client APIInterface
}

func (a *upstreamAPI) Create(upstream *adminv1.Upstream) (*adminv1.Upstream, *APIResponse) {
	out := &adminv1.Upstream{}
	err := a.client.Create(upstream, out)
	return out, err
}

func (a *upstreamAPI) Get(name string) (*adminv1.Upstream, *APIResponse) {
	out := &adminv1.Upstream{}
	err := a.client.Get(name, out)
	return out, err
}

func (a *upstreamAPI) List(params url.Values) (*adminv1.UpstreamList, error) {
	upstreamList := &adminv1.UpstreamList{}
	request := a.client.RestClient().Get().Resource("upstreams")
	for k, vals := range params {
		for _, v := range vals {
			request.Param(k, v)
		}
	}
	data, err := request.DoRaw()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, upstreamList); err != nil {
		return nil, err
	}

	if len(upstreamList.NextPage) > 0 {
		params.Add("offset", upstreamList.Offset)
		result, err := a.List(params)
		if err != nil {
			return nil, err
		}
		upstreamList.Items = append(upstreamList.Items, result.Items...)
	}

	return upstreamList, err
}

func (a *upstreamAPI) Delete(id string) error {
	return a.client.Delete(id)
}
