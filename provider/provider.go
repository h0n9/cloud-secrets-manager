package provider

type SecretProvider interface {
	Close() error
	ListSecrets(limit int) ([]string, error)
	GetSecretValue(secretID string) (string, error)
	SetSecretValue(secretID string, secretValue string) error
}
