package repositories

type SecretProviders interface {
	GetSecret(string) (string, error)
}
