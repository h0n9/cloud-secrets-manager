package cert

import (
	"fmt"

	"github.com/h0n9/toybox/cloud-secrets-manager/util"
	"github.com/spf13/cobra"
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
		err := util.GenerateAndSaveCertificate(service, namespace, certDir)
		if err != nil {
			return fmt.Errorf("faild to generate and save certificate")
		}
		fmt.Println("generated and saved certificates to " + certDir)
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
