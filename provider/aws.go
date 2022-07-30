package provider

import (
	"context"

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

func (provider *AWS) GetSecretValue(secretID string) (string, error) {
	req := &secretsmanager.GetSecretValueInput{SecretId: &secretID}
	resp, err := provider.client.GetSecretValue(provider.ctx, req)
	if err != nil {
		return "", err
	}
	return *resp.SecretString, nil
}
