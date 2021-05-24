package rootcmd

import (
	"context"

	"github.com/kong/kubernetes-ingress-controller/pkg/admission"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	"github.com/kong/kubernetes-ingress-controller/railgun/pkg/config"
)

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
			SecretGetter: &util.SecretGetterFromK8s{Reader: kubeclient},
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
