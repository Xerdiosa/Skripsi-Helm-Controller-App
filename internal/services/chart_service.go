package services

import (
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/repositories"
	"gorm.io/gorm"
)

type IChartService interface {
	InstallOrUpgradeChart(models.ChartRelease) error
	GetAllReleaseName() ([]string, error)
	GetReleaseDetail(string) (models.ChartRelease, error)
	RemoveChart(string) error
}

type ChartService struct {
	chartProvider repositories.Providers
}

func InitChartService(chartProvider repositories.Providers) IChartService {
	chartService := &ChartService{}
	chartService.chartProvider = chartProvider
	return chartService
}

func (h *ChartService) InstallOrUpgradeChart(chart models.ChartRelease) error {
	oldChartInterface, err := h.chartProvider.GetDetail(chart.ReleaseName)
	oldChart := oldChartInterface.(models.ChartRelease)
	if err == gorm.ErrRecordNotFound {
		return h.installChart(chart)
	}
	return h.upgradeChart(chart, oldChart)
}

func (h *ChartService) installChart(chart models.ChartRelease) error {
	chart.Revision = 1
	err := h.chartProvider.InstallComponent(chart)
	if err != nil {
		return err
	}
	err = h.chartProvider.Add(chart)
	return err
}

func (h *ChartService) upgradeChart(chart models.ChartRelease, oldChart models.ChartRelease) error {
	chart.Revision = oldChart.Revision + 1

	err := h.chartProvider.UpdateComponent(chart)
	if err != nil {
		return err
	}

	err = h.chartProvider.Update(chart)
	return err
}

func (h *ChartService) RemoveChart(chart string) error {
	chartInstance, err := h.chartProvider.GetDetail(chart)
	if err != nil {
		return err
	}
	err = h.chartProvider.UninstallComponent(chartInstance)
	if err != nil {
		return err
	}
	err = h.chartProvider.Remove(chartInstance)
	return err
}

func (h *ChartService) GetAllReleaseName() ([]string, error) {
	result, err := h.chartProvider.GetAllName()
	return result, err
}

func (h *ChartService) GetReleaseDetail(releaseName string) (models.ChartRelease, error) {
	resultInterface, err := h.chartProvider.GetDetail(releaseName)
	result := resultInterface.(models.ChartRelease)
	return result, err
}
