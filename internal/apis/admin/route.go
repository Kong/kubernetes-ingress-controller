package admin

import (
	"encoding/json"
	"net/url"

	adminv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/admin/v1"
)

type RouteGetter interface {
	Route() RouteInterface
}

type RouteInterface interface {
	List(url.Values) (*adminv1.RouteList, error)
	Get(string) (*adminv1.Route, *APIResponse)
	Create(*adminv1.Route) (*adminv1.Route, *APIResponse)
	Patch(string, *adminv1.Route) (*adminv1.Route, *APIResponse)
	Delete(string) error
}

type routeAPI struct {
	client APIInterface
}

func (a *routeAPI) Create(route *adminv1.Route) (*adminv1.Route, *APIResponse) {
	out := &adminv1.Route{}
	err := a.client.Create(route, out)
	return out, err
}

func (a *routeAPI) Patch(id string, route *adminv1.Route) (*adminv1.Route, *APIResponse) {
	out := &adminv1.Route{}
	err := a.client.Patch(id, route, out)
	return out, err
}

func (a *routeAPI) Get(name string) (*adminv1.Route, *APIResponse) {
	out := &adminv1.Route{}
	err := a.client.Get(name, out)
	return out, err
}

func (a *routeAPI) List(params url.Values) (*adminv1.RouteList, error) {
	if params == nil {
		params = url.Values{}
	}
	routeList := &adminv1.RouteList{}
	request := a.client.RestClient().Get().Resource("routes")
	for k, vals := range params {
		for _, v := range vals {
			request.Param(k, v)
		}
	}
	data, err := request.DoRaw()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, routeList); err != nil {
		return nil, err
	}

	if len(routeList.NextPage) > 0 {
		params.Set("offset", routeList.Offset)
		result, err := a.List(params)
		if err != nil {
			return nil, err
		}
		routeList.Items = append(routeList.Items, result.Items...)
	}

	return routeList, err
}

func (a *routeAPI) Delete(id string) error {
	return a.client.Delete(id)
}
