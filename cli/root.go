package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/h0n9/toybox/cloud-secrets-manager/cli/controller"
	"github.com/h0n9/toybox/cloud-secrets-manager/cli/injector"
	cliTemplate "github.com/h0n9/toybox/cloud-secrets-manager/cli/template"
)

var RootCmd = &cobra.Command{
	Use:   Name,
	Short: fmt.Sprintf("'%s' is a tool for playing with cloud-based secrets", Name),
}

func init() {
	cobra.EnableCommandSorting = false

	RootCmd.AddCommand(
		controller.Cmd,
		injector.Cmd,
		newLineCmd,
		cliTemplate.Cmd,
		newLineCmd,
		VersionCmd,
	)
}

var newLineCmd = &cobra.Command{Run: func(cmd *cobra.Command, args []string) {}} // new line
