package cli

import (
	"fmt"

	csm "github.com/h0n9/toybox/cloud-secrets-manager"
	"github.com/spf13/cobra"
)

const (
	DefaultProviderName   = "aws"
	DefaultTemplateBase64 = "e3sgcmFuZ2UgJGssICR2IDo9IC4gfX1be3sgJGsgfX1dCnt7ICR2IH19Cgp7eyBlbmQgfX0K"
	DefaultOutputFilename = "output"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("print '%s' version information", csm.Name),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(csm.Version)
	},
}
