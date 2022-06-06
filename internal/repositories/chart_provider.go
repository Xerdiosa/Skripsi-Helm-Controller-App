package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models/requests"
	helm "github.com/mittwald/go-helm-client"
	"gorm.io/gorm"
	"helm.sh/helm/v3/pkg/repo"
)

type ChartProvider struct {
	helmClient       map[string]helm.Client
	database         *gorm.DB
	defaultNamespace string
	chartRepo        repo.Entry
}

func InitChartProvider(helmClient map[string]helm.Client, database *gorm.DB, defaultNamespace string, chartRepo repo.Entry) Providers {
	chartProvider := &ChartProvider{}
	chartProvider.helmClient = helmClient
	chartProvider.database = database
	chartProvider.defaultNamespace = defaultNamespace
	chartProvider.chartRepo = chartRepo
	return chartProvider
}

func (h *ChartProvider) Convert(rawData interface{}) (interface{}, error) {
	jsonStr, err := json.Marshal(rawData)
	if err != nil {
		return nil, err
	}
	component := requests.ChartRelease{}
	err = json.Unmarshal(jsonStr, &component)
	if err != nil {
		return nil, err
	}
	if component.Namespace == "" {
		component.Namespace = h.defaultNamespace
	}
	return component.TransformToModels()
}

func (h *ChartProvider) PreProcess(data interface{}, prevData interface{}, module interface{}, moduleRelease interface{}) (interface{}, error) {
	processed, ok := data.(models.ChartRelease)
	if !ok {
		err := errors.New("conversion to chartRelease failed")
		return nil, err
	}
	releaseParsed, ok := moduleRelease.(models.ModuleRelease)
	if !ok {
		err := errors.New("conversion to release failed")
		return nil, err
	}
	var oldData models.ChartRelease
	if prevData != nil {
		oldData, ok = prevData.(models.ChartRelease)
		if !ok {
			err := errors.New("conversion to chartRelease failed")
			return nil, err
		}
	}

	processed.ModuleReleaseID = releaseParsed.ID

	processed.Revision = oldData.Revision + 1
	return processed, nil
}

func (h *ChartProvider) InstallComponent(chartInterface interface{}) error {
	chart, ok := chartInterface.(models.ChartRelease)
	if !ok {
		err := errors.New("conversion to chartRelease failed")
		return err
	}

	chartSpec := helm.ChartSpec{
		ReleaseName: chart.ReleaseName,
		ChartName:   chart.Name,
		Version:     chart.Version,
		UpgradeCRDs: true,
		Wait:        true,
		Timeout:     time.Minute * 5,
		ValuesYaml:  chart.Values,
		Namespace:   chart.Namespace,
	}

	if _, ok := h.helmClient[chart.Namespace]; !ok {
		err := errors.New("unknown namespace")
		return err
	}

	if err := h.helmClient[chart.Namespace].AddOrUpdateChartRepo(h.chartRepo); err != nil {
		return err
	}

	_, err := h.helmClient[chart.Namespace].InstallOrUpgradeChart(context.Background(), &chartSpec)
	return err
}

func (h *ChartProvider) UpdateComponent(releaseInterface interface{}) error {
	return h.InstallComponent(releaseInterface)
}

func (h *ChartProvider) UninstallComponent(releaseInterface interface{}) error {
	release, ok := releaseInterface.(models.ChartRelease)
	if !ok {
		err := errors.New("conversion to chartRelease failed")
		return err
	}

	if _, ok := h.helmClient[release.Namespace]; !ok {
		err := errors.New("unknown namespace")
		return err
	}
	err := h.helmClient[release.Namespace].UninstallReleaseByName(release.ReleaseName)
	return err
}

func (h *ChartProvider) GetAllName() ([]string, error) {
	var names []string
	result := h.database.Model(&models.ChartRelease{}).Pluck("release_name", &names)
	if result.Error != nil {
		return nil, result.Error
	}
	return names, nil
}

func (h *ChartProvider) Add(releaseInterface interface{}) error {
	release, ok := releaseInterface.(models.ChartRelease)
	if !ok {
		err := errors.New("conversion to chartRelease failed")
		return err
	}

	result := h.database.Create(&release)
	return result.Error
}

func (h *ChartProvider) Remove(releaseInterface interface{}) error {
	release, ok := releaseInterface.(models.ChartRelease)
	if !ok {
		err := errors.New("conversion to chartRelease failed")
		return err
	}

	result := h.database.Delete(&models.ChartRelease{}, "release_name = ?", release.ReleaseName)
	return result.Error
}

func (h *ChartProvider) Update(releaseInterface interface{}) error {
	release, ok := releaseInterface.(models.ChartRelease)
	if !ok {
		err := errors.New("conversion to chartRelease failed")
		return err
	}

	result := h.database.Model(&release).Where("release_name = ?", release.ReleaseName).Updates(release)
	return result.Error
}

func (h *ChartProvider) GetDetailFromComponent(releaseInterface interface{}) (interface{}, error) {
	release, ok := releaseInterface.(models.ChartRelease)
	if !ok {
		err := errors.New("conversion to chartRelease failed")
		return nil, err
	}

	return h.GetDetail(release.ReleaseName)
}

func (h *ChartProvider) GetDetail(releaseName string) (interface{}, error) {
	var release models.ChartRelease
	result := h.database.Where("release_name = ?", releaseName).First(&release)
	return release, result.Error
}

func (h *ChartProvider) GetFromModuleReleaseID(ModuleReleaseID uint) ([]interface{}, error) {
	var charts []models.ChartRelease
	result := h.database.Where("module_release_id = ?", ModuleReleaseID).Find(&charts)

	var chartsInterface []interface{} = make([]interface{}, len(charts))
	for i, v := range charts {
		chartsInterface[i] = v
	}

	return chartsInterface, result.Error
}
