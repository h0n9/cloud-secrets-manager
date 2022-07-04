package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	Name    = "cloud-secrets-manager"
	Version = "v0.0.4"

	DefaultProviderName   = "aws"
	DefaultTemplateBase64 = "e3sgcmFuZ2UgJGssICR2IDo9IC4gfX1be3sgJGsgfX1dCnt7ICR2IH19Cgp7eyBlbmQgfX0K"
	DefaultOutputFilename = "output"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("print '%s' version information", Name),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
