package provider

type SecretProvider interface {
	GetSecretValue(string) (string, error)
}
