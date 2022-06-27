package controller

import (
	"os"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	csiWebhook "github.com/h0n9/toybox/cloud-secrets-manager/webhook"
)

var Cmd = &cobra.Command{
	Use:   "controller",
	Short: "admission webhook controller",
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a server for managing admission webhooks",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLogger(zap.New())
		logger := log.Log.WithName("cloud-secrets-injector-controller")

		mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{Logger: logger})
		if err != nil {
			logger.Error(err, "faild to setup controller")
			os.Exit(1)
		}

		hookServer := mgr.GetWebhookServer()
		hookServer.Register("/mutate", &webhook.Admission{Handler: &csiWebhook.Mutator{
			Client: mgr.GetClient(),
		}})
		hookServer.Register("/validate", &webhook.Admission{Handler: &csiWebhook.Validator{
			Client: mgr.GetClient(),
		}})
		logger.Info("registered mutate, validator handlers to /mutate, /validate webhook uris")

		logger.Info("starting controller")
		err = mgr.Start(signals.SetupSignalHandler())
		if err != nil {
			logger.Error(err, "failed to run controller")
			os.Exit(1)
		}
	},
}

func init() {
	Cmd.AddCommand(runCmd)
}
