package provider

type SecretProvider interface {
	Close() error
	ListSecrets() ([]string, error)
	GetSecretValue(string) (string, error)
	SetSecretValue(string, string) error
}
