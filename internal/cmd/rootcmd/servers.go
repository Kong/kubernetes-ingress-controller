package rootcmd

import (
	"context"

	"github.com/bombsimon/logrusr"
	"github.com/kong/deck/file"

	"github.com/kong/kubernetes-ingress-controller/internal/admission"
	"github.com/kong/kubernetes-ingress-controller/internal/diagnostics"
	"github.com/kong/kubernetes-ingress-controller/internal/manager"
	"github.com/kong/kubernetes-ingress-controller/internal/util"
)

func StartAdmissionServer(ctx context.Context, c *manager.Config) error {
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

func StartDiagnosticsServer(ctx context.Context, port int, c *manager.Config) (diagnostics.Server, error) {
	deprecatedLogger, err := util.MakeLogger(c.LogLevel, c.LogFormat)
	if err != nil {
		return diagnostics.Server{}, err
	}
	logger := logrusr.NewLogger(deprecatedLogger)

	if !c.EnableProfiling && !c.EnableConfigDumps {
		logger.Info("diagnostics server disabled")
		return diagnostics.Server{}, nil
	}

	s := diagnostics.Server{
		Logger:           logger,
		ProfilingEnabled: c.EnableProfiling,
	}
	if c.EnableConfigDumps {
		s.ConfigDumps = util.ConfigDumpDiagnostic{
			DumpsIncludeSensitive: c.DumpSensitiveConfig,
			SuccessfulConfigs:     make(chan file.Content),
			FailedConfigs:         make(chan file.Content),
		}
	}
	go func() {
		if err := s.Listen(ctx, port); err != nil {
			logger.Error(err, "unable to start diagnostics server")
		}
	}()
	return s, nil
}
