package cert

import (
	"context"
	"fmt"
	"path"

	"github.com/spf13/cobra"
	apimetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	admissionregistrationv1 "k8s.io/client-go/applyconfigurations/admissionregistration/v1"
	metav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	csm "github.com/h0n9/toybox/cloud-secrets-manager"
	"github.com/h0n9/toybox/cloud-secrets-manager/util"
)

var Cmd = &cobra.Command{
	Use:   "cert",
	Short: "certificate manager",
}

var (
	namespace string
	service   string
	certDir   string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate CA, server certificates with private key",
	RunE: func(cmd *cobra.Command, args []string) error {
		// init
		ctx := context.Background()
		kubeCli := kubernetes.NewForConfigOrDie(config.GetConfigOrDie())

		// generate self-signed CA certificate
		fmt.Printf("generating certificates ... ")
		caCert, serverCert, serverPrivKey, err := util.GenerateCertificate(service, namespace)
		if err != nil {
			fmt.Println("❌")
			return err
		}
		fmt.Println("✅")

		// write certificates to files
		for filename, data := range map[string][]byte{
			"ca.crt":  caCert,
			"tls.crt": serverCert,
			"tls.key": serverPrivKey,
		} {
			filename = path.Join(certDir, filename)
			fmt.Printf("writing to '%s' ... ", filename)
			err = util.WriteBytesToFile(filename, data)
			if err != nil {
				fmt.Println("❌")
				return err
			}
			fmt.Println("✅")
		}

		fmt.Printf("applying MutatingWebhookConfiguration '%s'... ", service)
		_, err = kubeCli.AdmissionregistrationV1().
			MutatingWebhookConfigurations().
			Apply(ctx, admissionregistrationv1.
				MutatingWebhookConfiguration(service).
				WithWebhooks(
					admissionregistrationv1.MutatingWebhook().
						WithName(csm.AnnotationPrefix).
						WithClientConfig(admissionregistrationv1.WebhookClientConfig().
							WithCABundle(caCert...).
							WithService(admissionregistrationv1.ServiceReference().
								WithNamespace(namespace).
								WithName(service).
								WithPath("/mutate").
								WithPort(8443),
							),
						).
						WithSideEffects("None").
						WithAdmissionReviewVersions("v1beta1").
						WithFailurePolicy("Fail").
						WithRules(admissionregistrationv1.RuleWithOperations().
							WithAPIGroups("").
							WithAPIVersions("v1").
							WithOperations("CREATE", "UPDATE").
							WithResources("pods").
							WithScope("Namespaced"),
						).
						WithNamespaceSelector(metav1.LabelSelector().
							WithMatchLabels(map[string]string{
								"cloud-secrets-injector": "enabled",
							}),
						),
				),
				apimetav1.ApplyOptions{
					FieldManager: csm.Name,
				},
			)
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
		&certDir,
		"cert-dir",
		"/etc/certs",
		"directory containing certificate and private key files",
	)
	Cmd.AddCommand(generateCmd)
}
