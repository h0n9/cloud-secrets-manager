package cli

import (
	"fmt"

	csm "github.com/h0n9/cloud-secrets-manager"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("print '%s' version information", csm.Name),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(csm.Version)
	},
}
