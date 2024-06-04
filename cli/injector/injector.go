package injector

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/h0n9/cloud-secrets-manager/handler"
	"github.com/h0n9/cloud-secrets-manager/provider"
	"github.com/h0n9/cloud-secrets-manager/util"
)

const (
	DefaultProviderName   = "aws"
	DefaultTemplateBase64 = "e3sgcmFuZ2UgJGssICR2IDo9IC4gfX1be3sgJGsgfX1dCnt7ICR2IH19Cgp7eyBlbmQgfX0K"
	DefaultOutputFilename = "output"
)

var Cmd = &cobra.Command{
	Use:   "injector",
	Short: "cloud-based secrets injector",
}

var (
	providerName              string
	secretID                  string
	templateBase64            string
	output                    string
	decodeBase64EncodedSecret bool
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "inject cloud-based secrets into containers",
	RunE: func(cmd *cobra.Command, args []string) error {
		// init logger
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

		// init context
		ctx := context.Background()

		logger.Info().Msg("initialized context")

		// get envs
		if secretID == "" {
			return fmt.Errorf("failed to read 'secret-id' flag")
		}

		logger.Info().Msg("read environment variables")

		// decode base64-encoded template to string
		templateStr, err := util.DecodeBase64(templateBase64)
		if err != nil {
			return err
		}
		tmpl := template.New("secret-template")
		tmpl, err = tmpl.Parse(templateStr)
		if err != nil {
			return err
		}

		logger.Info().Msg("loaded template")

		var (
			secretProvider provider.SecretProvider
			secretHandler  *handler.SecretHandler
		)

		switch strings.ToLower(providerName) {
		case "aws":
			secretProvider, err = provider.NewAWS(ctx)
		case "gcp":
			secretProvider, err = provider.NewGCP(ctx)
		default:
			return fmt.Errorf("failed to figure out secret provider")
		}
		if err != nil {
			return err
		}
		defer secretProvider.Close()

		logger.Info().Msg(fmt.Sprintf("initialized secret provider '%s'", providerName))

		secretHandler, err = handler.NewSecretHandler(secretProvider, tmpl)
		if err != nil {
			return err
		}

		logger.Info().Msg("initialized secret handler")

		err = secretHandler.Save(secretID, output, decodeBase64EncodedSecret)
		if err != nil {
			return err
		}

		logger.Info().Msg(fmt.Sprintf("saved secret id '%s' values to '%s'", secretID, output))

		return nil
	},
}

func init() {
	runCmd.Flags().StringVar(
		&providerName,
		"provider",
		"aws",
		"cloud provider name",
	)
	runCmd.Flags().StringVar(
		&secretID,
		"secret-id",
		"",
		"secret id",
	)
	runCmd.Flags().StringVar(
		&templateBase64,
		"template",
		"e3sgcmFuZ2UgJGssICR2IDo9IC4gfX1be3sgJGsgfX1dCnt7ICR2IH19Cgp7eyBlbmQgfX0K",
		"base64 encoded template string",
	)
	runCmd.Flags().StringVar(
		&output,
		"output",
		"secret",
		"output filename",
	)
	runCmd.Flags().BoolVar(
		&decodeBase64EncodedSecret,
		"decode-b64-secret",
		false,
		"decode base64 encoded secret",
	)
	Cmd.AddCommand(runCmd)
}
