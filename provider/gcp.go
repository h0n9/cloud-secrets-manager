package provider

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type GCP struct {
	ctx    context.Context
	client *secretmanager.Client
}

func NewGCP(ctx context.Context) (*GCP, error) {
	// create new client with default options
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GCP{
		ctx:    ctx,
		client: client,
	}, nil
}

func (provider *GCP) Close() error {
	return provider.client.Close()
}

// The secretID in the format `projects/*/secrets/*/versions/*`.
// `projects/*/secrets/*/versions/latest`: recently created
func (provider *GCP) GetSecretValue(secretID string) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{Name: secretID}
	resp, err := provider.client.AccessSecretVersion(provider.ctx, req)
	if err != nil {
		return "", err
	}
	return string(resp.GetPayload().GetData()), nil
}

func (provider *GCP) SetSecretValue(secretID, secretValue string) error {
	return nil
}
