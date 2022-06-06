package repositories

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

type VaultSecretProvider struct {
	vaultClient *api.Client
}

func InitVaultSecretProvider(vaultClient *api.Client) SecretProviders {
	vaultSecretProvider := &VaultSecretProvider{}
	vaultSecretProvider.vaultClient = vaultClient
	return vaultSecretProvider
}

func (v VaultSecretProvider) GetSecret(key string) (string, error) {
	vaultKV := strings.SplitN(key, ":", 2)
	if len(vaultKV) != 2 {
		return "", fmt.Errorf("key %s parsing error", key)
	}
	vaultSecret, err := v.vaultClient.Logical().Read(vaultKV[0])
	if err != nil {
		return "nil", err
	}
	data := vaultSecret.Data["data"].(map[string]interface{})
	value, ok := data[vaultKV[1]]
	if !ok {
		return "", fmt.Errorf("key %s does not exist", vaultKV[1])
	}
	return value.(string), nil
}
