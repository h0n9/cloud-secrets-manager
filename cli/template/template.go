package template

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "template",
	Short: "template related operations",
}

func init() {
	Cmd.AddCommand(
		encodeCmd,
		decodeCmd,
		testCmd,
	)
}
