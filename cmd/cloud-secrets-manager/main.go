package main

import (
	"os"

	"github.com/h0n9/cloud-secrets-manager/cli"
)

func main() {
	err := cli.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
