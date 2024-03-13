package provider

type SecretProvider interface {
	Close() error
	GetSecretValue(string) (string, error)
	SetSecretValue(string) error
}
