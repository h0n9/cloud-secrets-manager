package secrets

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "secrets",
	Short: "CLI for managing secrets",
}

var (
	providerName string
	secretID     string
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit a secret",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	editCmd.Flags().StringVar(
		&providerName,
		"provider",
		"aws",
		"cloud provider name",
	)
	editCmd.Flags().StringVar(
		&secretID,
		"secret-id",
		"",
		"secret id",
	)
	Cmd.AddCommand(editCmd)
}
