package cert

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/h0n9/toybox/cloud-secrets-manager/util"
)

var Cmd = &cobra.Command{
	Use:   "cert",
	Short: "certificate manager",
}

var (
	service   string
	namespace string
	certDir   string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate CA, server certificates with private key",
	RunE: func(cmd *cobra.Command, args []string) error {
		// generate self-signed CA certificate
		fmt.Printf("generating and saving certificates to %s ... ", certDir)
		caCertPEM, err := util.GenerateAndSaveCertificate(service, namespace, certDir)
		if err != nil {
			fmt.Println("❌")
			return err
		}
		fmt.Println("✅")

		kubeCli := kubernetes.NewForConfigOrDie(config.GetConfigOrDie())
		ctx := context.Background()

		fmt.Printf("getting MutatingWebhookConfiguration ... ")
		mutatingWebhookConfiguration, err := kubeCli.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(ctx, service, metav1.GetOptions{})
		if err != nil {
			fmt.Println("❌")
			return err
		}
		fmt.Println("✅")

		// mutate CABundle
		for i, _ := range mutatingWebhookConfiguration.Webhooks {
			mutatingWebhookConfiguration.Webhooks[i].ClientConfig.CABundle = caCertPEM
		}

		fmt.Printf("updating MutatingWebhookConfiguration ... ")
		_, err = kubeCli.AdmissionregistrationV1().MutatingWebhookConfigurations().Update(ctx, mutatingWebhookConfiguration, metav1.UpdateOptions{})
		if err != nil {
			fmt.Println("❌")
			return err
		}
		fmt.Println("✅")

		return nil
	},
}

func init() {
	generateCmd.Flags().StringVar(
		&service,
		"service",
		"cloud-secrets-controller",
		"kubernetes service resource's name",
	)
	generateCmd.Flags().StringVar(
		&namespace,
		"namespace",
		"cloud-secrets-controller",
		"kubernetes service resource's namespace",
	)
	generateCmd.Flags().StringVar(
		&certDir,
		"cert-dir",
		"/etc/certs",
		"directory containing certificate and private key files",
	)
	Cmd.AddCommand(generateCmd)
}
