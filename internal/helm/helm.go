package helm

import (
	"flag"
	"path/filepath"

	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/gudangada/data-warehouse/warehouse-controller/internal/configs"
	helm "github.com/mittwald/go-helm-client"
)

var helmClient helm.Client

func GetHelmClient(authConfig configs.AuthConfig) (helm.Client, error) {
	// func GetHelmClient() (helm.Client, error) {

	if helmClient == nil {
		var config *rest.Config
		var err error

		switch authConfig.Method {
		// switch "kubeconfig" {

		case "kubeconfig":
			var kubeconfig *string
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
			} else {
				kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
			}
			flag.Parse()

			config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
			if err != nil {
				return nil, err
			}
		case "service-account":
			config, err = rest.InClusterConfig()
			if err != nil {
				panic(err.Error())
			}
		}

		opt := &helm.RestConfClientOptions{
			Options:    &helm.Options{},
			RestConfig: config,
		}

		helmClient, err = helm.NewClientFromRestConf(opt)
		if err != nil {
			panic(err)
		}

		gudangadaBiRepo := repo.Entry{
			Name: "gudangada-bi",
			URL:  "s3://gudangada-bi-helm-charts/stable",
		}

		if err := helmClient.AddOrUpdateChartRepo(gudangadaBiRepo); err != nil {
			panic(err)
		}
	}
	return helmClient, nil
}

func GenerateHelmClient(authConfig configs.AuthConfig) (map[string]helm.Client, error) {
	helmClient := map[string]helm.Client{}
	var config *rest.Config
	var err error

	switch authConfig.Method {
	case "kubeconfig":
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
	case "service-account":
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	for _, namespace := range authConfig.AvailableNamespace {
		opt := &helm.RestConfClientOptions{
			Options: &helm.Options{
				Namespace: namespace,
			},
			RestConfig: config,
		}

		helmClientTmp, err := helm.NewClientFromRestConf(opt)
		if err != nil {
			panic(err)
		}

		helmClient[namespace] = helmClientTmp
	}
	return helmClient, nil
}

func GetChartRepo(chartRepoConfig configs.ChartRepo) repo.Entry {
	chartRepo := repo.Entry{
		Name: chartRepoConfig.Name,
		URL:  chartRepoConfig.URL,
	}
	return chartRepo
}
