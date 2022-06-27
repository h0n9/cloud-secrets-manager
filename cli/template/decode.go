package template

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/h0n9/toybox/cloud-secrets-manager/util"
)

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "decode base64-encoded template string to human-readable string",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] == "" {
			return fmt.Errorf("failed to decode empty string")
		}
		template, err := util.DecodeBase64(args[0])
		if err != nil {
			return err
		}
		fmt.Println(template)
		return nil
	},
}
