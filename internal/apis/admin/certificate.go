package admin

import (
	"encoding/json"
	"net/url"

	adminv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/admin/v1"
)

type CertificateGetter interface {
	Certificate() CertificateInterface
}

type CertificateInterface interface {
	List(url.Values) (*adminv1.CertificateList, error)
	Get(string) (*adminv1.Certificate, *APIResponse)
	Create(*adminv1.Certificate) (*adminv1.Certificate, *APIResponse)
	Patch(string, *adminv1.Certificate) (*adminv1.Certificate, *APIResponse)
	Delete(string) error
}

type certificateAPI struct {
	client APIInterface
}

func (a *certificateAPI) Create(target *adminv1.Certificate) (*adminv1.Certificate, *APIResponse) {
	out := &adminv1.Certificate{}
	err := a.client.Put(target.ID, target, out)
	return out, err
}

func (a *certificateAPI) Get(name string) (*adminv1.Certificate, *APIResponse) {
	out := &adminv1.Certificate{}
	err := a.client.Get(name, out)
	return out, err
}

func (a *certificateAPI) Patch(id string, cert *adminv1.Certificate) (*adminv1.Certificate, *APIResponse) {
	out := &adminv1.Certificate{}
	err := a.client.Patch(id, cert, out)
	return out, err
}

func (a *certificateAPI) List(params url.Values) (*adminv1.CertificateList, error) {
	if params == nil {
		params = url.Values{}
	}

	list := &adminv1.CertificateList{}
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
	if err := json.Unmarshal(data, list); err != nil {
		return nil, err
	}

	if len(list.NextPage) > 0 {
		params.Set("offset", list.Offset)
		result, err := a.List(params)
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, result.Items...)
	}

	return list, err
}

func (a *certificateAPI) Delete(name string) error {
	return a.client.Delete(name)
}
