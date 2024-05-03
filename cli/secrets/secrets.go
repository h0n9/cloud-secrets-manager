package secrets

import (
	"github.com/spf13/cobra"
)

const (
	DefaultProviderName = "aws"
	DefaultEditor       = "vim"
)

var (
	providerName string
)

var Cmd = &cobra.Command{
	Use:   "secrets",
	Short: "CLI for managing secrets",
}

func init() {
	Cmd.AddCommand(listCmd)
	Cmd.AddCommand(editCmd)
}
