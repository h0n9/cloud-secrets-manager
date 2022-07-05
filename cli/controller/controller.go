package controller

import (
	"context"
	"net/http"
	"os"
	"path"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/admissionregistration/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/h0n9/toybox/cloud-secrets-manager/util"
	csiWebhook "github.com/h0n9/toybox/cloud-secrets-manager/webhook"
)

var Cmd = &cobra.Command{
	Use:   "controller",
	Short: "admission webhook controller",
}

var (
	service       string
	namespace     string
	port          int
	certDir       string
	injectorImage string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a server for managing admission webhooks",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLogger(zap.New())
		logger := log.Log.WithName(service)

		mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
			Logger:                 logger,
			HealthProbeBindAddress: ":8081",
			Port:                   port,
			CertDir:                certDir,
		})
		if err != nil {
			logger.Error(err, "faild to setup controller")
			os.Exit(1)
		}
		kubeCli := mgr.GetClient()
		ctx := context.Background()

		// read ca.crt
		caCertPEM, err := util.ReadFileToStr(path.Join(certDir, "ca.crt"))
		if err != nil {
			logger.Error(err, "failed to read 'ca.crt'")
		}
		_ = util.EncodeBase64(caCertPEM)

		// get current mutatingWebhookConfiguration
		mutatingWebhookConfiguration := v1.MutatingWebhookConfiguration{}
		err = kubeCli.Get(ctx, client.ObjectKey{
			Namespace: namespace,
			Name:      service,
		}, &mutatingWebhookConfiguration)
		if err != nil {
			logger.Error(err, "failed to get mutatingWebhookConfiguration")
		}

		err = mgr.AddHealthzCheck("healthz", func(req *http.Request) error { return nil })
		if err != nil {
			logger.Error(err, "failed to add healthz check")
		}
		err = mgr.AddReadyzCheck("readyz", func(req *http.Request) error { return nil })
		if err != nil {
			logger.Error(err, "failed to add readyz check")
		}
		logger.Info("added healthz, readyz probes")

		hookServer := mgr.GetWebhookServer()
		hookServer.Register("/mutate", &webhook.Admission{Handler: &csiWebhook.Mutator{
			Client:        mgr.GetClient(),
			InjectorImage: injectorImage,
		}})
		hookServer.Register("/validate", &webhook.Admission{Handler: &csiWebhook.Validator{
			Client: mgr.GetClient(),
		}})

		logger.Info("starting controller")
		err = mgr.Start(signals.SetupSignalHandler())
		if err != nil {
			logger.Error(err, "failed to run controller")
			os.Exit(1)
		}
	},
}

func init() {
	runCmd.Flags().StringVar(
		&service,
		"service",
		"cloud-secrets-controller",
		"kubernetes service resource's name",
	)
	runCmd.Flags().StringVar(
		&namespace,
		"namespace",
		"cloud-secrets-controller",
		"kubernetes service resource's namespace",
	)
	runCmd.Flags().IntVar(
		&port,
		"port",
		8443,
		"port for webhook server to listen on",
	)
	runCmd.Flags().StringVar(
		&certDir,
		"cert-dir",
		"/etc/certs",
		"directory containing certificate and private key files",
	)
	runCmd.Flags().StringVar(
		&injectorImage,
		"image",
		"ghcr.io/h0n9/cloud-secrets-manager:v0.0.4",
		"docker image name with tag for init container",
	)
	Cmd.AddCommand(runCmd)
}
