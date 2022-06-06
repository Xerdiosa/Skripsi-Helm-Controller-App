package configs

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfigs struct {
	Server     ServerConfig `yaml:"server"`
	Database   DBConfig     `yaml:"database"`
	Kubernetes AuthConfig   `yaml:"kubernetes"`
	Vault      VaultConfig  `yaml:"vault"`
	ChartRepo  ChartRepo    `yaml:"chartRepo"`
}

type ServerConfig struct {
	Host string `yaml:"host" env:"SERVER_HOST" env-default:"127.0.0.1"`
	Port int    `yaml:"port" env:"SERVER_PORT" env-default:"3000"`
}

type DBConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_NAME"`
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"pass" env:"DB_PASS"`
}

type AuthConfig struct {
	Method             string   `yaml:"method" env:"KUBERNETES_AUTHENTICATION_METHOD" env-default:"kubeconfig"`
	DefaultNamespace   string   `yaml:"defaultNamespace" env:"KUBERNETES_DEFAULT_NAMESPACE" env-devault:"default"`
	AvailableNamespace []string `yaml:"availableNamespace" env:"KUBERNETES_AVAILABLE_NAMESPACE" env-devault:"default"`
}

type VaultConfig struct {
	URL            string              `yaml:"url" env:"VAULT_URL"`
	AuthMethod     string              `yaml:"authMethod" env:"VAULT_AUTH_METHOD"`
	KubeMethod     VaultKubeMethod     `yaml:"kubeMethod"`
	UserPassMethod VaultUserPassMethod `yaml:"userPassMethod"`
}

type VaultKubeMethod struct {
	RoleName  string `yaml:"roleName" env:"VAULT_KUBE_ROLE_NAME"`
	MountPath string `yaml:"mountPath" env:"VAULT_MOUNT_PATH"`
}

type VaultUserPassMethod struct {
	Username string `yaml:"username" env:"VAULT_USERNAME"`
	Password string `yaml:"password" env:"VAULT_PASSWORD"`
}

type ChartRepo struct {
	Name string `yaml:"name" env:"CHART_REPO_NAME"`
	URL  string `yaml:"url" env:"CHART_REPO_URL"`
}

func InitAppConfigs() (*AppConfigs, error) {
	var appConfigs AppConfigs

	err := cleanenv.ReadEnv(&appConfigs)
	if err != nil {
		return nil, err
	}

	configFile := "config.yaml"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &appConfigs, nil
	}

	err = cleanenv.ReadConfig(configFile, &appConfigs)
	if err != nil {
		return nil, err
	}

	return &appConfigs, nil
}
