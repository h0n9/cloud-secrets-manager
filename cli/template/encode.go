package template

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/h0n9/cloud-secrets-manager/util"
)

var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "encode human-readable template string to base64-encoded string",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] == "" {
			return fmt.Errorf("failed to encode empty string")
		}
		fmt.Println(util.EncodeBase64StrToStr(args[0]))
		return nil
	},
}
