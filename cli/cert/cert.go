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
	namespace   string
	serviceName string
	secretName  string
	certDir     string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate CA, server certificates with private key",
	RunE: func(cmd *cobra.Command, args []string) error {
		// init
		ctx := context.Background()
		kubeCli := kubernetes.NewForConfigOrDie(config.GetConfigOrDie())

		// get secret for certificates
		fmt.Printf("getting Secret '%s'... ", secretName)
		secret, err := kubeCli.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
		if err != nil {
			fmt.Println("❌")
			return err
		}
		fmt.Println("✅")

		var (
			caCertPEM        []byte = secret.Data["ca.crt"]
			serverCertPEM    []byte = secret.Data["tls.crt"]
			serverPrivKeyPEM []byte = secret.Data["tls.key"]
		)

		generate := len(caCertPEM) < 10 || len(serverCertPEM) < 10 || len(serverPrivKeyPEM) < 10
		if generate {
			// generate self-signed CA certificate
			fmt.Printf("generating certificates ... ")
			caCertPEM, serverCertPEM, serverPrivKeyPEM, err = util.GenerateCertificate(secretName, namespace, certDir)
			if err != nil {
				fmt.Println("❌")
				return err
			}
			fmt.Println("✅")

			secret.Data["ca.crt"] = caCertPEM
			secret.Data["tls.crt"] = serverCertPEM
			secret.Data["tls.key"] = serverPrivKeyPEM

			// update secret for certificates
			fmt.Printf("updating Secret '%s'... ", secretName)
			_, err = kubeCli.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
			if err != nil {
				fmt.Println("❌")
				return err
			}
			fmt.Println("✅")
		}

		fmt.Printf("getting MutatingWebhookConfiguration '%s' ... ", serviceName)
		mutatingWebhookConfiguration, err := kubeCli.AdmissionregistrationV1().MutatingWebhookConfigurations().Get(ctx, serviceName, metav1.GetOptions{})
		if err != nil {
			fmt.Println("❌")
			return err
		}
		fmt.Println("✅")

		// mutate CABundle
		for i, _ := range mutatingWebhookConfiguration.Webhooks {
			mutatingWebhookConfiguration.Webhooks[i].ClientConfig.CABundle = caCertPEM
		}

		fmt.Printf("updating MutatingWebhookConfiguration '%s' ... ", serviceName)
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
		"cloud-secrets-controller",
		"kubernetes service resource's namespace",
	)
	generateCmd.Flags().StringVar(
		&serviceName,
		"service-name",
		"cloud-secrets-controller",
		"kubernetes service resource's name",
	)
	generateCmd.Flags().StringVar(
		&secretName,
		"secret-name",
		"cloud-secrets-controller-tls",
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
