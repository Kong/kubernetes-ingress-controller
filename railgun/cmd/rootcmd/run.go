package rootcmd

import (
	"context"
	"fmt"

	"github.com/kong/kubernetes-ingress-controller/pkg/admission"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	"github.com/kong/kubernetes-ingress-controller/railgun/manager"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/config"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Run(ctx context.Context, c *config.Config) error {
	if err := StartAdmissionServer(ctx, c); err != nil {
		return fmt.Errorf("StartAdmissionServer: %w", err)
	}
	return manager.Run(ctx, c)
}

func StartAdmissionServer(ctx context.Context, c *config.Config) error {
	log, err := util.MakeLogger(c.LogLevel, c.LogFormat)
	if err != nil {
		return err
	}
	if c.AdmissionServer.ListenAddr == "off" {
		log.Info("admission webhook server disabled")
		return nil
	}
	kubeclient, err := c.GetKubeClient()
	if err != nil {
		return err
	}
	kongclient, err := c.GetKongClient(ctx)
	if err != nil {
		return err
	}
	srv, err := admission.MakeTLSServer(&c.AdmissionServer, &admission.RequestHandler{
		Validator: admission.KongHTTPValidator{
			ConsumerSvc:  kongclient.Consumers,
			PluginSvc:    kongclient.Plugins,
			Logger:       log,
			SecretGetter: &secretGetter{Reader: kubeclient},
		},
	})
	if err != nil {
		return err
	}
	go func() {
		err := srv.ListenAndServeTLS("", "")
		log.WithError(err).Error("admission webhook server stopped")
	}()
	return nil
}

type secretGetter struct {
	Reader client.Reader
}

func (s *secretGetter) GetSecret(namespace string, name string) (*corev1.Secret, error) {
	var res corev1.Secret
	err := s.Reader.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, &res)
	return &res, err
}
