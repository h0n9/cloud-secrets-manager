package provider

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AWS struct {
	ctx    context.Context
	cfg    config.Config
	client *secretsmanager.Client
}

func NewAWS(ctx context.Context) (*AWS, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := secretsmanager.NewFromConfig(cfg)
	return &AWS{
		ctx:    ctx,
		cfg:    cfg,
		client: client,
	}, nil
}

func (provider *AWS) Close() error {
	return nil
}

func (provider *AWS) ListSecrets() ([]string, error) {
	req := &secretsmanager.ListSecretsInput{}
	var secrets []string

	for {
		// list secrets
		resp, err := provider.client.ListSecrets(provider.ctx, req)
		if err != nil {
			return nil, err
		}

		// append secret names
		for _, secret := range resp.SecretList {
			secrets = append(secrets, *secret.Name)
		}

		// break if no more secrets
		if resp.NextToken == nil {
			break
		}

		// set next token
		req.NextToken = resp.NextToken
	}

	return secrets, nil
}

func (provider *AWS) GetSecretValue(secretID string) (string, error) {
	req := &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretID)}
	resp, err := provider.client.GetSecretValue(provider.ctx, req)
	if err != nil {
		return "", err
	}
	return *resp.SecretString, nil
}

func (provider *AWS) SetSecretValue(secretID, secretValue string) error {
	req := &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(secretID),
		SecretString: aws.String(secretValue),
	}
	_, err := provider.client.PutSecretValue(provider.ctx, req)
	if err != nil {
		return err
	}
	return nil
}
