package admin

import (
	"encoding/json"
	"net/url"

	adminv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/admin/v1"
)

type ConsumerGetter interface {
	Consumers() ConsumerInterface
}

type ConsumerInterface interface {
	List(url.Values) (*adminv1.ConsumerList, error)
	Get(string) (*adminv1.Consumer, *APIResponse)
	Create(*adminv1.Consumer) (*adminv1.Consumer, *APIResponse)
	Patch(string, *adminv1.Consumer) (*adminv1.Consumer, *APIResponse)
	Delete(string) error
}

type consumerAPI struct {
	client APIInterface
}

func (a *consumerAPI) Create(consumer *adminv1.Consumer) (*adminv1.Consumer, *APIResponse) {
	out := &adminv1.Consumer{}
	err := a.client.Put(consumer.ID, consumer, out)
	return out, err
}

func (a *consumerAPI) Patch(id string, consumer *adminv1.Consumer) (*adminv1.Consumer, *APIResponse) {
	out := &adminv1.Consumer{}
	err := a.client.Patch(id, consumer, out)
	return out, err
}

func (a *consumerAPI) Get(name string) (*adminv1.Consumer, *APIResponse) {
	out := &adminv1.Consumer{}
	err := a.client.Get(name, out)
	return out, err
}

func (a *consumerAPI) List(params url.Values) (*adminv1.ConsumerList, error) {
	if params == nil {
		params = url.Values{}
	}

	ConsumerList := &adminv1.ConsumerList{}
	request := a.client.RestClient().Get().Resource("consumers")
	for k, vals := range params {
		for _, v := range vals {
			request.Param(k, v)
		}
	}
	data, err := request.DoRaw()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, ConsumerList); err != nil {
		return nil, err
	}

	if len(ConsumerList.NextPage) > 0 {
		params.Set("offset", ConsumerList.Offset)
		result, err := a.List(params)
		if err != nil {
			return nil, err
		}
		ConsumerList.Items = append(ConsumerList.Items, result.Items...)
	}

	return ConsumerList, err
}

func (a *consumerAPI) Delete(id string) error {
	return a.client.Delete(id)
}
