package secrets

import (
	"context"
	"crypto/sha1"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/h0n9/cloud-secrets-manager/provider"
	"github.com/h0n9/cloud-secrets-manager/util"
)

const (
	DefaultProviderName = "aws"
	DefaultEditor       = "vim"
)

var Cmd = &cobra.Command{
	Use:   "secrets",
	Short: "CLI for managing secrets",
}

var (
	providerName string
	secretID     string
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit a secret",
	RunE: func(cmd *cobra.Command, args []string) error {
		// check secretID
		if secretID == "" {
			return fmt.Errorf("failed to read 'secret-id' flag")
		}

		// define variables
		var (
			err            error
			secretProvider provider.SecretProvider
		)

		// init ctx
		ctx := context.Background()

		// init secret provider
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

		// get secret value
		secretValue, err := secretProvider.GetSecretValue(secretID)
		if err != nil {
			return err
		}

		// convert json to yaml
		data, err := yaml.JSONToYAML([]byte(secretValue))
		if err != nil {
			return err
		}

		// write data to tmp file
		UserCacheDir, err := os.UserCacheDir()
		if err != nil {
			return err
		}
		hash := sha1.Sum([]byte(secretID))
		tmpFilePath := path.Join(UserCacheDir, fmt.Sprintf("%x", hash))
		err = os.WriteFile(tmpFilePath, data, 0644)
		if err != nil {
			return err
		}

		// get initial stat of tmp file
		initialStat, err := os.Stat(tmpFilePath)
		if err != nil {
			return err
		}

		// open tmp file with editor(e.g. vim)
		editor := util.GetEnv("EDITOR", DefaultEditor)
		execCmd := exec.Command(editor, tmpFilePath)
		execCmd.Stdin = os.Stdin
		execCmd.Stdout = os.Stdout
		err = execCmd.Run()
		if err != nil {
			return err
		}

		// get updated stat of tmp file
		updatedStat, err := os.Stat(tmpFilePath)
		if err != nil {
			return err
		}

		// check if tmp file is updated. if not, return nil
		if initialStat.ModTime().Equal(updatedStat.ModTime()) &&
			initialStat.Size() == updatedStat.Size() {
			fmt.Println("found nothing to update")
			return nil
		}

		// read tmp file
		data, err = os.ReadFile(tmpFilePath)
		if err != nil {
			return err
		}

		// convert yaml to json
		data, err = yaml.YAMLToJSON(data)
		if err != nil {
			return err
		}

		// update secret value
		secretValue = string(data)

		// set secret value to provider
		err = secretProvider.SetSecretValue(secretID, secretValue)
		if err != nil {
			return err
		}

		// remove tmp file
		err = os.Remove(tmpFilePath)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	editCmd.Flags().StringVar(
		&providerName,
		"provider",
		DefaultProviderName,
		"cloud provider name",
	)
	editCmd.Flags().StringVar(
		&secretID,
		"secret-id",
		"",
		"secret id",
	)
	Cmd.AddCommand(editCmd)
}
