package secrets

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/h0n9/cloud-secrets-manager/provider"
)

const (
	DefaultListSecretsLimit = 100
)

var (
	listSecretsLimit int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		// define variables
		var (
			err            error
			secretProvider provider.SecretProvider
		)

		// init ctx
		ctx := context.Background()

		// init secret provider
		switch strings.ToLower(providerName) {
		case "aws":
			secretProvider, err = provider.NewAWS(ctx)
		case "gcp":
			secretProvider, err = provider.NewGCP(ctx)
		default:
			return fmt.Errorf("failed to figure out secret provider")
		}
		if err != nil {
			return err
		}
		defer secretProvider.Close()

		// list secrets
		secrets, err := secretProvider.ListSecrets(listSecretsLimit)
		if err != nil {
			return err
		}

		// print secrets
		for _, secret := range secrets {
			fmt.Println(secret)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringVar(
		&providerName,
		"provider",
		DefaultProviderName,
		"cloud provider name",
	)
	listCmd.Flags().IntVar(
		&listSecretsLimit,
		"limit",
		DefaultListSecretsLimit,
		"limit the number of secrets to list",
	)
}
