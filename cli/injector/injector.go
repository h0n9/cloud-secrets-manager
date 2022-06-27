package injector

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/h0n9/toybox/cloud-secrets-manager/handler"
	"github.com/h0n9/toybox/cloud-secrets-manager/provider"
	"github.com/h0n9/toybox/cloud-secrets-manager/util"
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
		secretId := util.GetEnv("SECRET_ID", "")
		if secretId == "" {
			return fmt.Errorf("failed to read 'SECRET_ID'")
		}
		providerName := util.GetEnv("PROVIDER_NAME", DefaultProviderName)
		templateBase64 := util.GetEnv("TEMPLATE_BASE64", DefaultTemplateBase64)
		templateFile := util.GetEnv("TEMPLATE_FILE", "")
		outputFile := util.GetEnv("OUTPUT_FILE", DefaultOutputFilename)

		logger.Info().Msg("read environment variables")

		// decode base64-encoded template to string
		templateStr, err := util.DecodeBase64(templateBase64)
		if err != nil {
			return err
		}
		if templateFile != "" {
			templateStr, err = util.ReadFileToStr(templateFile)
			if err != nil {
				return err
			}
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
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("failed to figure out secret provider")
		}

		logger.Info().Msg(fmt.Sprintf("initialized secret provider '%s'", providerName))

		secretHandler, err = handler.NewSecretHandler(secretProvider, tmpl)
		if err != nil {
			return err
		}

		logger.Info().Msg("initialized secret handler")

		err = secretHandler.Save(secretId, outputFile)
		if err != nil {
			return err
		}

		logger.Info().Msg(fmt.Sprintf("saved secret id '%s' values to '%s'", secretId, outputFile))

		return nil
	},
}

func init() {
	Cmd.AddCommand(runCmd)
}
