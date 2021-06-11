package rootcmd

import (
	"context"

	"github.com/bombsimon/logrusr"
	"github.com/go-logr/logr"
	"github.com/kong/kubernetes-ingress-controller/pkg/admission"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
	"github.com/kong/kubernetes-ingress-controller/railgun/internal/diagnostics"
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

func StartProfilingServer(ctx context.Context, c *config.Config) error {
	deprecatedLogger, err := util.MakeLogger(c.LogLevel, c.LogFormat)

	var logger logr.Logger = logrusr.NewLogger(deprecatedLogger)
	if err != nil {
		return err
	}
	if c.EnableProfiling {
		s := diagnostics.Server{Logger: logger}
		go func() {
			if err := s.Listen(ctx); err != nil {
				logger.Error(err, "unable to start diagnostics server")
			}
		}()
	}
	return nil
}
