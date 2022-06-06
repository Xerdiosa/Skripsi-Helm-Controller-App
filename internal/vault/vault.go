package vault

import (
	"context"
	"fmt"

	"github.com/gudangada/data-warehouse/warehouse-controller/internal/configs"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/kubernetes"
	"github.com/hashicorp/vault/api/auth/userpass"
)

var vaultClient *api.Client

func GetVaultSecret(vaultConfig configs.VaultConfig) (*api.Client, error) {
	if vaultClient != nil {
		return vaultClient, nil
	}
	config := api.DefaultConfig()

	config.Address = vaultConfig.URL

	vaultClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	var auth api.AuthMethod

	switch vaultConfig.AuthMethod {
	case "service-account":
		auth, err = kubernetes.NewKubernetesAuth(
			vaultConfig.KubeMethod.RoleName,
			kubernetes.WithMountPath(vaultConfig.KubeMethod.MountPath),
		)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize Kubernetes auth method: %w", err)
		}
	case "userpass":
		auth, err = userpass.NewUserpassAuth(
			vaultConfig.UserPassMethod.Username,
			&userpass.Password{FromString: vaultConfig.UserPassMethod.Password},
		)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize UserPass auth method: %w", err)
		}
	}

	authInfo, err := vaultClient.Auth().Login(context.TODO(), auth)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with auth: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	return vaultClient, nil
}
