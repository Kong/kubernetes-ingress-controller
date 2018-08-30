package admin

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"

	adminv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/admin/v1"
)

type PluginGetter interface {
	Plugins() PluginInterface
}

type PluginInterface interface {
	List(url.Values) (*adminv1.PluginList, error)
	Get(string) (*adminv1.Plugin, *APIResponse)
	CreateInRoute(string, *adminv1.Plugin) (*adminv1.Plugin, *APIResponse)
	CreateInService(string, *adminv1.Plugin) (*adminv1.Plugin, *APIResponse)
	CreateGlobal(*adminv1.Plugin) (*adminv1.Plugin, *APIResponse)
	Patch(string, *adminv1.Plugin) (*adminv1.Plugin, *APIResponse)
	Delete(string) error

	GetByID(string) (*adminv1.Plugin, error)

	GetAllByRoute(string) ([]adminv1.Plugin, error)
	GetByRoute(string, string) (*adminv1.Plugin, error)

	GetAllByService(serviceID string) ([]adminv1.Plugin, error)
	GetByService(name, serviceID string) (*adminv1.Plugin, error)
	GetAllByServiceWitConsumer(serviceID string) ([]*adminv1.Plugin, error)
}

type pluginAPI struct {
	client APIInterface
}

func (a *pluginAPI) create(name, id string, plugin *adminv1.Plugin) (*adminv1.Plugin, *APIResponse) {
	out := &adminv1.Plugin{}
	rawData, err := json.Marshal(plugin)
	if err != nil {
		return nil, &APIResponse{err: err}
	}
	resp := a.client.
		RestClient().
		Post().
		Resource(name).
		SubResource(id, "plugins").
		Body(rawData).
		SetHeader("Content-Type", "application/json").
		Do()

	statusCode := reflect.ValueOf(resp).FieldByName("statusCode").Int()
	raw, err := resp.Raw()
	response := &APIResponse{StatusCode: int(statusCode), err: err}
	if err != nil {
		response.Raw = raw
	}

	response.err = json.Unmarshal(raw, out)
	return out, response
}

func (a *pluginAPI) CreateInRoute(id string, plugin *adminv1.Plugin) (*adminv1.Plugin, *APIResponse) {
	return a.create("routes", id, plugin)
}

func (a *pluginAPI) CreateInService(id string, plugin *adminv1.Plugin) (*adminv1.Plugin, *APIResponse) {
	return a.create("services", id, plugin)
}

func (a *pluginAPI) CreateGlobal(plugin *adminv1.Plugin) (*adminv1.Plugin, *APIResponse) {
	out := &adminv1.Plugin{}
	err := a.client.Create(plugin, out)
	return out, err
}

func (a *pluginAPI) Patch(id string, route *adminv1.Plugin) (*adminv1.Plugin, *APIResponse) {
	out := &adminv1.Plugin{}
	err := a.client.Patch(id, route, out)
	return out, err
}

func (a *pluginAPI) Get(name string) (*adminv1.Plugin, *APIResponse) {
	out := &adminv1.Plugin{}
	err := a.client.Get(name, out)
	return out, err
}

func (a *pluginAPI) List(params url.Values) (*adminv1.PluginList, error) {
	if params == nil {
		params = url.Values{}
	}

	PluginList := &adminv1.PluginList{}
	request := a.client.RestClient().Get().Resource("plugins")
	for k, vals := range params {
		for _, v := range vals {
			request.Param(k, v)
		}
	}
	data, err := request.DoRaw()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, PluginList); err != nil {
		return nil, err
	}

	if len(PluginList.NextPage) > 0 {
		params.Set("offset", PluginList.Offset)
		result, err := a.List(params)
		if err != nil {
			return nil, err
		}
		PluginList.Items = append(PluginList.Items, result.Items...)
	}

	return PluginList, err
}

func (a *pluginAPI) Delete(id string) error {
	return a.client.Delete(id)
}

func (a *pluginAPI) GetByID(id string) (*adminv1.Plugin, error) {
	plugins, err := a.List(nil)
	if err != nil {
		return nil, err
	}

	for _, plugin := range plugins.Items {
		if plugin.ID == id {
			return &plugin, nil
		}
	}

	return nil, &PluginNotConfiguredError{
		Message: fmt.Sprintf("Plugin %v is not configured", id),
	}
}

func (a *pluginAPI) listByResource(id, resource string, params url.Values) (*adminv1.PluginList, error) {
	plugins := &adminv1.PluginList{}
	request := a.client.RestClient().
		Get().
		Resource(resource).
		SubResource(id, "plugins")

	for k, vals := range params {
		for _, v := range vals {
			request.Param(k, v)
		}
	}
	data, err := request.DoRaw()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, plugins); err != nil {
		return nil, err
	}

	if len(plugins.NextPage) > 0 {
		params.Add("offset", plugins.Offset)
		result, err := a.listByResource(id, resource, params)
		if err != nil {
			return nil, err
		}
		plugins.Items = append(plugins.Items, result.Items...)
	}

	return plugins, nil
}

func (a *pluginAPI) GetByRoute(name, routeID string) (*adminv1.Plugin, error) {
	plugins, err := a.listByResource(routeID, "routes", nil)
	if err != nil {
		return nil, err
	}

	for _, plugin := range plugins.Items {
		if plugin.Name == name {
			return &plugin, nil
		}
	}

	return nil, &PluginNotConfiguredError{
		Message: fmt.Sprintf("Plugin %v is not configured in the route %v", name, routeID),
	}
}

func (a *pluginAPI) GetAllByRoute(routeID string) ([]adminv1.Plugin, error) {
	list, err := a.listByResource(routeID, "routes", nil)
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

func (a *pluginAPI) GetByService(name, serviceID string) (*adminv1.Plugin, error) {
	plugins, err := a.listByResource(serviceID, "services", nil)
	if err != nil {
		return nil, err
	}

	for _, plugin := range plugins.Items {
		if plugin.Name == name &&
			plugin.Service == serviceID {
			return &plugin, nil
		}
	}

	return nil, &PluginNotConfiguredError{
		Message: fmt.Sprintf("Plugin %v is not configured in the service %v", name, serviceID),
	}
}

func (a *pluginAPI) GetAllByService(serviceID string) ([]adminv1.Plugin, error) {
	plugins, err := a.listByResource(serviceID, "services", nil)
	if err != nil {
		return nil, err
	}

	return plugins.Items, nil
}

func (a *pluginAPI) GetAllByServiceWitConsumer(serviceID string) ([]*adminv1.Plugin, error) {
	plugins, err := a.listByResource(serviceID, "services", nil)
	if err != nil {
		return nil, err
	}

	res := make([]*adminv1.Plugin, 0)
	for _, plugin := range plugins.Items {
		if plugin.Service == serviceID && plugin.Consumer != "" {
			res = append(res, &plugin)
		}
	}

	return res, nil
}

// PluginNotConfiguredError defines an
type PluginNotConfiguredError struct {
	Message string
}

// check to verify PluginNotConfiguredError implements the Error interface
var _ error = &PluginNotConfiguredError{}

// Error implements the Error interface.
func (e PluginNotConfiguredError) Error() string {
	return e.Message
}

// IsPluginNotConfiguredError checks if the type of the error is PluginNotConfiguredError
func IsPluginNotConfiguredError(err error) bool {
	_, ok := err.(*PluginNotConfiguredError)
	return ok
}
