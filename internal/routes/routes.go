package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/configs"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/controllers"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/database"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/helm"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/kinesis"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/repositories"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/services"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/vault"
)

type Route struct{}

func (r *Route) Init(config configs.AppConfigs) *mux.Router {
	helmClient, err := helm.GenerateHelmClient(config.Kubernetes)
	if err != nil {
		panic(err)
	}

	chartRepo := helm.GetChartRepo(config.ChartRepo)

	kinesisClient, err := kinesis.GetKinesisClient()
	if err != nil {
		panic(err)
	}

	database, err := database.GetDB(config.Database)
	if err != nil {
		panic(err)
	}

	vault, err := vault.GetVaultSecret(config.Vault)
	if err != nil {
		panic(err)
	}

	defaultNamespace := config.Kubernetes.DefaultNamespace

	moduleRepository := repositories.InitModuleRepository(database)

	chartProvider := repositories.InitChartProvider(helmClient, database, defaultNamespace, chartRepo)
	kinesisProvider := repositories.InitKinesisProvider(database, kinesisClient)

	vaultSecretProvider := repositories.InitVaultSecretProvider(vault)

	componentProviders := map[string]repositories.Providers{
		"chart":   chartProvider,
		"kinesis": kinesisProvider,
	}

	secretProviders := map[string]repositories.SecretProviders{
		"vault": vaultSecretProvider,
	}

	chartService := services.InitChartService(chartProvider)
	kinesisService := services.InitKinesisService(kinesisProvider)
	moduleService := services.InitModuleService(moduleRepository, componentProviders, secretProviders)

	chartController := controllers.InitChartController(chartService)
	kinesisController := controllers.InitKinesisController(kinesisService)
	moduleController := controllers.InitModuleController(moduleService)

	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/chart", chartController.Release).Methods(http.MethodPost)
	router.HandleFunc("/chart", chartController.GetAllReleaseName).Methods(http.MethodGet)
	router.HandleFunc("/chart/{chart-name}", chartController.GetReleaseDetail).Methods(http.MethodGet)
	router.HandleFunc("/chart/{chart-name}", chartController.UpdateRelease).Methods(http.MethodPut)
	router.HandleFunc("/chart/{chart-name}", chartController.RemoveRelease).Methods(http.MethodDelete)

	router.HandleFunc("/kinesis", kinesisController.Release).Methods(http.MethodPost)
	router.HandleFunc("/kinesis", kinesisController.GetAllReleaseName).Methods(http.MethodGet)
	router.HandleFunc("/kinesis/{kinesis-name}", kinesisController.GetReleaseDetail).Methods(http.MethodGet)
	router.HandleFunc("/kinesis/{kinesis-name}", kinesisController.UpdateRelease).Methods(http.MethodPut)
	router.HandleFunc("/kinesis/{kinesis-name}", kinesisController.RemoveRelease).Methods(http.MethodDelete)

	router.HandleFunc("/module", moduleController.AddModule).Methods(http.MethodPost)
	router.HandleFunc("/module/release", moduleController.AddModuleRelease).Methods(http.MethodPost)
	router.HandleFunc("/module/release", moduleController.GetAllReleaseName).Methods(http.MethodGet)
	router.HandleFunc("/module/release/{release-name}", moduleController.GetReleaseDetail).Methods(http.MethodGet)
	router.HandleFunc("/module/release/{release-name}", moduleController.UpdateModuleRelease).Methods(http.MethodPut)
	router.HandleFunc("/module/release/{release-name}", moduleController.DeleteModuleRelease).Methods(http.MethodDelete)

	return router
}
