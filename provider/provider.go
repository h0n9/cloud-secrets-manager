package provider

type SecretProvider interface {
	Close() error
	ListSecrets(int) ([]string, error)
	GetSecretValue(string) (string, error)
	SetSecretValue(string, string) error
}
