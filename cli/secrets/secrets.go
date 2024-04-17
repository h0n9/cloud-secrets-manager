package secrets

import (
	"github.com/spf13/cobra"
)

const (
	DefaultProviderName = "aws"
	DefaultEditor       = "vim"
)

var Cmd = &cobra.Command{
	Use:   "secrets",
	Short: "CLI for managing secrets",
}

func init() {
	Cmd.AddCommand(editCmd)
}
