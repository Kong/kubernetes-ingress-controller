package controller

import (
	"reflect"
	"testing"

	"github.com/hbagdi/go-kong/kong"
	configurationv1 "github.com/kong/kubernetes-ingress-controller/internal/apis/configuration/v1"
)

func TestMergeRouteAndKongIngress(t *testing.T) {
	r := &kong.Route{
		Hosts: []*string{kong.String("foo.com"), kong.String("bar.com")},
		Paths: []*string{kong.String("/")},
	}

	kongIngress := &configurationv1.KongIngress{
		Route: &configurationv1.Route{},
	}
	updated := mergeRouteAndKongIngress(r, kongIngress)
	if updated {
		t.Errorf("expected: false, actual: %v", updated)
	}

	kongIngress.Route.RegexPriority = 10
	kongIngress.Route.PreserveHost = true
	updated = mergeRouteAndKongIngress(r, kongIngress)
	if !updated {
		t.Errorf("expected: false, actual: %v", updated)
	}
	if *r.RegexPriority != 10 {
		t.Errorf("expected regex priority to be 10")
	}
	if !*r.PreserveHost {
		t.Errorf("expected PreserveHost to be true")
	}

	kongIngress.Route.Protocols = []string{"https"}
	updated = mergeRouteAndKongIngress(r, kongIngress)
	if !updated {
		t.Errorf("expected: false, actual: %v", updated)
	}
	if len(r.Protocols) != 1 {
		t.Errorf("expected length to be 1")
	}
	if *r.Protocols[0] != "https" {
		t.Errorf("expected protocols to be 'https'")
	}

	updated = mergeRouteAndKongIngress(r, kongIngress)
	if updated {
		t.Errorf("expected updated to be false on no-op")
	}

	kongIngress.Route.Methods = []string{"GET", "POST"}
	updated = mergeRouteAndKongIngress(r, kongIngress)
	if !updated {
		t.Errorf("expected updated to be true")
	}

	if !reflect.DeepEqual(toStringArray(r.Methods), kongIngress.Route.Methods) {
		t.Errorf("expected %v, got %v", kongIngress.Route.Methods, toStringArray(r.Methods))
	}
}
