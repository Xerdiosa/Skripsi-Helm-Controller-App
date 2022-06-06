package repositories

import (
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models"
	"gorm.io/gorm"
)

type IModuleRepository interface {
	InsertModule(models.Module) error
	InsertModuleRelease(models.ModuleRelease) (models.ModuleRelease, error)
	GetModule(string, string) (models.Module, error)
	GetModuleRelease(string) (models.ModuleRelease, error)
	GetAllModuleRelease() ([]string, error)
	DeleteModuleRelease(models.ModuleRelease) error
}

type ModuleRepository struct {
	database *gorm.DB
}

func InitModuleRepository(database *gorm.DB) IModuleRepository {
	ModuleRepository := &ModuleRepository{}
	ModuleRepository.database = database
	return ModuleRepository
}

func (m ModuleRepository) InsertModule(module models.Module) error {
	result := m.database.Create(&module)
	return result.Error
}

func (m ModuleRepository) InsertModuleRelease(moduleRelease models.ModuleRelease) (models.ModuleRelease, error) {
	result := m.database.Create(&moduleRelease)
	return moduleRelease, result.Error
}

func (m ModuleRepository) GetModule(moduleName string, version string) (models.Module, error) {
	var module models.Module
	var result *gorm.DB
	if version != "" {
		result = m.database.Where("name = ? AND version = ?", moduleName, version).First(&module)
	} else {
		result = m.database.Order("created_at desc").Where("name = ?", moduleName, version).First(&module)

	}
	return module, result.Error
}

func (m ModuleRepository) GetModuleRelease(moduleReleaseName string) (models.ModuleRelease, error) {
	var moduleRelease models.ModuleRelease
	result := m.database.Order("created_at desc").Where("name = ?", moduleReleaseName).First(&moduleRelease)
	return moduleRelease, result.Error
}

func (m ModuleRepository) GetAllModuleRelease() ([]string, error) {
	var names []string
	result := m.database.Model(&models.ModuleRelease{}).Pluck("name", &names)
	if result.Error != nil {
		return nil, result.Error
	}
	return names, nil
}

func (m ModuleRepository) DeleteModuleRelease(moduleRelease models.ModuleRelease) error {
	result := m.database.Delete(&moduleRelease)
	return result.Error
}
