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
	namespace string
	service   string
	secret    string
	certDir   string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate CA, server certificates with private key",
	RunE: func(cmd *cobra.Command, args []string) error {
		// init
		ctx := context.Background()
		kubeCli := kubernetes.NewForConfigOrDie(config.GetConfigOrDie())

		// get secret for certificates
		fmt.Printf("getting Secret '%s'... ", secret)
		kubeSecret, err := kubeCli.CoreV1().Secrets(namespace).Get(ctx, secret, metav1.GetOptions{})
		if err != nil {
			fmt.Println("❌")
			return err
		}
		fmt.Println("✅")

		var (
			caCertPEM        []byte = kubeSecret.Data["ca.crt"]
			serverCertPEM    []byte = kubeSecret.Data["tls.crt"]
			serverPrivKeyPEM []byte = kubeSecret.Data["tls.key"]
		)

		generate := len(caCertPEM) < 10 || len(serverCertPEM) < 10 || len(serverPrivKeyPEM) < 10
		if generate {
			// generate self-signed CA certificate
			fmt.Printf("generating certificates ... ")
			caCertPEM, serverCertPEM, serverPrivKeyPEM, err = util.GenerateCertificate(service, namespace, certDir)
			if err != nil {
				fmt.Println("❌")
				return err
			}
			fmt.Println("✅")

			kubeSecret.Data["ca.crt"] = caCertPEM
			kubeSecret.Data["tls.crt"] = serverCertPEM
			kubeSecret.Data["tls.key"] = serverPrivKeyPEM

			// update secret for certificates
			fmt.Printf("updating Secret '%s'... ", secret)
			_, err = kubeCli.CoreV1().Secrets(namespace).Update(ctx, kubeSecret, metav1.UpdateOptions{})
			if err != nil {
				fmt.Println("❌")
				return err
			}
			fmt.Println("✅")
		}

		fmt.Printf("getting MutatingWebhookConfiguration '%s' ... ", service)
		mutatingWebhookConfiguration, err := kubeCli.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(ctx, service, metav1.GetOptions{})
		if err != nil {
			fmt.Println("❌")
			return err
		}
		fmt.Println("✅")

		// mutate CABundle
		for i := range mutatingWebhookConfiguration.Webhooks {
			mutatingWebhookConfiguration.Webhooks[i].ClientConfig.CABundle = caCertPEM
		}

		fmt.Printf("updating MutatingWebhookConfiguration '%s' ... ", service)
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
		&namespace,
		"namespace",
		"cloud-secrets-manager",
		"kubernetes service resource's namespace",
	)
	generateCmd.Flags().StringVar(
		&service,
		"service",
		"cloud-secrets-manager",
		"kubernetes service resource's name",
	)
	generateCmd.Flags().StringVar(
		&secret,
		"secret",
		"cloud-secrets-manager-tls",
		"kubernetes secret resource's name",
	)
	generateCmd.Flags().StringVar(
		&certDir,
		"cert-dir",
		"/etc/certs",
		"directory containing certificate and private key files",
	)
	Cmd.AddCommand(generateCmd)
}
