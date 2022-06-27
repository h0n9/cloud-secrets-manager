package template

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/h0n9/toybox/cloud-secrets-manager/util"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "test base64-encoded template string",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] == "" {
			return fmt.Errorf("failed to test empty string")
		}
		templateStr, err := util.DecodeBase64(args[0])
		if err != nil {
			return err
		}
		tmpl := template.New("sample-template-to-test")
		tmpl, err = tmpl.Parse(templateStr)
		if err != nil {
			return err
		}
		err = tmpl.Execute(os.Stdout, map[string]string{
			"hello":      "world",
			"life":       "is beautiful",
			"difference": "makes our life more rich",
		})
		if err != nil {
			return err
		}
		return nil
	},
}
