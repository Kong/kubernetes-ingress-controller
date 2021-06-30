package seeder

import (
	"context"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	knative "knative.dev/networking/pkg/client/clientset/versioned"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kong/kubernetes-ingress-controller/railgun/internal/proxy"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/clientset"
)

// -----------------------------------------------------------------------------
// Seeder
// -----------------------------------------------------------------------------

// Seeder is an object which can perform pre-fetch seed rounds to cache
// supported objects into the proxy cache server.
type Seeder struct {
	namespaces       []string
	ingressClassName string

	logger logrus.FieldLogger
	prx    proxy.Proxy

	kc       *kubernetes.Clientset
	kongc    *clientset.Clientset
	knativec *knative.Clientset
}

// New provides a need *Seeder object.
func New(restCFG *rest.Config, prx proxy.Proxy) (*Seeder, error) {
	return NewBuilder(restCFG, prx).Build()
}

// Seed lists all supported API types, filters them to make sure the
// object is supported (e.g. uses the managers' ingress.class) and then
// pulls fresh copies of the object from the Kubernetes API to (re)seed
// the proxy cache. This is required to deal with situations where events
// get lost by controllers due to networking failures, poorly timed controller
// pod restarts, e.t.c.
func (s *Seeder) Seed(ctx context.Context) error {
	// FIXME - optionality/enablement for apis
	objs := make([]client.Object, 0)

	s.logger.Info("fetching supported core kubernetes objects")
	coreObjs, err := s.fetchCore(ctx)
	if err != nil {
		return err
	}
	objs = append(objs, coreObjs...)

	s.logger.Info("fetching supported kong kubernetes objects")
	kongObjs, err := s.fetchKong(ctx)
	if err != nil {
		return err
	}
	objs = append(objs, kongObjs...)

	s.logger.Info("fetching 3rd party kubernetes objects")
	otherObjs, err := s.fetchOther(ctx)
	if err != nil {
		return err
	}
	objs = append(objs, otherObjs...)

	if len(objs) < 1 {
		s.logger.Info("seed round successful: there were no objects that needed to be cached")
		return nil
	}

	if err := s.prx.UpdateObjects(objs...); err != nil {
		return err
	}
	s.logger.Infof("seed round successful: %d objects were added to the cache", len(objs))

	return nil
}

// -----------------------------------------------------------------------------
// Seeder - Controller Runtime - Runnable Implementation
// -----------------------------------------------------------------------------

func (s *Seeder) Start(ctx context.Context) error {
	return s.Seed(ctx)
}
