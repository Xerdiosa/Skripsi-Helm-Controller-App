package services

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/repositories"
	"sigs.k8s.io/yaml"
)

type IModuleService interface {
	InstallModule(models.Module) error
	ReleaseModule(models.Module, models.ModuleRelease, bool) error
	UpdateModuleRelease(models.Module, models.ModuleRelease, bool) error
	DeleteModuleRelease(models.ModuleRelease) error
	GetAllReleaseName() ([]string, error)
	GetReleaseDetail(releaseName string) (models.ModuleRelease, error)
}

type ModuleService struct {
	moduleRepository repositories.IModuleRepository
	providers        map[string]repositories.Providers
	secretProviders  map[string]repositories.SecretProviders
}

func InitModuleService(moduleRepository repositories.IModuleRepository, providers map[string]repositories.Providers, secretProviders map[string]repositories.SecretProviders) IModuleService {
	moduleService := &ModuleService{}
	moduleService.moduleRepository = moduleRepository
	moduleService.providers = providers
	moduleService.secretProviders = secretProviders
	return moduleService
}

func (m ModuleService) InstallModule(module models.Module) error {

	err := m.moduleRepository.InsertModule(module)
	if err != nil {
		return err
	}
	return nil
}

func (m ModuleService) ReleaseModule(module models.Module, release models.ModuleRelease, deleteOnFail bool) error {
	module, err := m.moduleRepository.GetModule(module.Name, module.Version)
	if err != nil {
		return err
	}

	release.ModuleID = module.ID
	release.ModuleName = module.Name
	release.Revision = 1

	finalSpec, err := m.applyChartTemplate(module, module.Spec, release)
	if err != nil {
		return err
	}

	var spec map[string][]interface{}
	err = yaml.Unmarshal([]byte(finalSpec), &spec)
	if err != nil {
		return err
	}

	release, err = m.moduleRepository.InsertModuleRelease(release)
	if err != nil {
		return err
	}

	for handler := range spec {
		if _, ok := m.providers[handler]; !ok {
			err := errors.New("component handler not implemented")
			return err
		}
	}

	for handler, components := range spec {
		for i, component := range components {
			spec[handler][i], err = m.providers[handler].Convert(component)
			if err != nil {
				return err
			}
		}
	}

	for handler, components := range spec {
		for i := range components {
			spec[handler][i], err = m.providers[handler].PreProcess(spec[handler][i], nil, module, release)
			if err != nil {
				return err
			}
		}
	}

	for handler, components := range spec {
		for _, component := range components {
			err = m.providers[handler].InstallComponent(component)
			if err != nil {
				if deleteOnFail {
					m.providers[handler].UninstallComponent(component)
					m.forceDelete(release)
				}
				return err
			}
			err = m.providers[handler].Add(component)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m ModuleService) UpdateModuleRelease(module models.Module, release models.ModuleRelease, deleteOnFail bool) error {
	module, err := m.moduleRepository.GetModule(module.Name, module.Version)
	if err != nil {
		return err
	}

	oldRelease, err := m.moduleRepository.GetModuleRelease(release.Name)
	if err != nil {
		return err
	}

	release.ModuleID = module.ID
	release.ModuleName = module.Name
	release.Revision = oldRelease.Revision + 1

	finalSpec, err := m.applyChartTemplate(module, module.Spec, release)
	if err != nil {
		return err
	}

	var spec map[string][]interface{}
	err = yaml.Unmarshal([]byte(finalSpec), &spec)
	if err != nil {
		return err
	}

	release, err = m.moduleRepository.InsertModuleRelease(release)
	if err != nil {
		return err
	}

	err = m.moduleRepository.DeleteModuleRelease(oldRelease)
	if err != nil {
		return err
	}

	for handler := range spec {
		if _, ok := m.providers[handler]; !ok {
			err := errors.New("component handler not implemented")
			return err
		}
	}

	for handler, components := range spec {
		for i, component := range components {
			spec[handler][i], err = m.providers[handler].Convert(component)
			if err != nil {
				return err
			}
		}
	}

	for handler, components := range spec {
		for i := range components {
			oldChart, err := m.providers[handler].GetDetailFromComponent(spec[handler][i])
			if err != nil {
				return err
			}

			spec[handler][i], err = m.providers[handler].PreProcess(spec[handler][i], oldChart, module, release)
			if err != nil {
				return err
			}
		}
	}

	for handler, components := range spec {
		for _, component := range components {
			err = m.providers[handler].UpdateComponent(component)
			if err != nil {
				if deleteOnFail {
					m.forceDelete(release)
				}
				return err
			}
			err = m.providers[handler].Update(component)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *ModuleService) GetAllReleaseName() ([]string, error) {
	return h.moduleRepository.GetAllModuleRelease()
}

func (h *ModuleService) GetReleaseDetail(releaseName string) (models.ModuleRelease, error) {
	return h.moduleRepository.GetModuleRelease(releaseName)
}

func (h *ModuleService) applyChartTemplate(chart models.Module, chartTemplate string, release models.ModuleRelease) (string, error) {
	templateVal := models.ModuleTemplate{
		Module:  chart.Name,
		Version: chart.Version,
		Release: release.Name,
	}
	err := yaml.Unmarshal([]byte(release.Values), &templateVal.Values)
	if err != nil {
		return "", err
	}

	if secret, ok := templateVal.Values["secret"]; ok {
		parsedSecret := make(map[string]map[string]interface{})
		mappedSecret := secret.(map[string]interface{})
		for secretProviderName, rawSecret := range mappedSecret {
			if _, ok := h.secretProviders[secretProviderName]; !ok {
				err := errors.New("component handler not implemented")
				return "", err
			}

			secretProvider := h.secretProviders[secretProviderName]
			parsedRawSecret := rawSecret.(map[string]interface{})
			parsedSecret[secretProviderName], err = h.getSecret(parsedRawSecret, secretProvider)
			if err != nil {
				return "", err
			}

			templateVal.Values["secret"] = parsedSecret
		}
	}

	buf := new(bytes.Buffer)
	tmpl, err := template.New("template").Funcs(funcMap()).Parse(chartTemplate)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(buf, templateVal)
	if err != nil {
		return "", err
	}
	applied := buf.String()
	return applied, nil
}

func (m ModuleService) getSecret(secretList map[string]interface{}, secretProvider repositories.SecretProviders) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for key, element := range secretList {
		secret, err := secretProvider.GetSecret(element.(string))
		if err != nil {
			return nil, err
		}
		result[key] = secret
	}
	return result, nil
}

func (m ModuleService) forceDelete(release models.ModuleRelease) {
	release, err := m.moduleRepository.GetModuleRelease(release.Name)
	if err != nil {
		return
	}
	for _, handler := range m.providers {
		components, _ := handler.GetFromModuleReleaseID(release.ID)
		for _, component := range components {
			handler.UninstallComponent(component)
			handler.Remove(component)
		}
	}

	m.moduleRepository.DeleteModuleRelease(release)
}

func (m ModuleService) DeleteModuleRelease(release models.ModuleRelease) error {
	release, err := m.moduleRepository.GetModuleRelease(release.Name)
	if err != nil {
		return err
	}

	for _, handler := range m.providers {
		components, err := handler.GetFromModuleReleaseID(release.ID)
		if err != nil {
			return err
		}
		for _, component := range components {
			err = handler.UninstallComponent(component)
			if err != nil {
				return err
			}

			err = handler.Remove(component)
			if err != nil {
				return err
			}
		}
	}

	err = m.moduleRepository.DeleteModuleRelease(release)
	if err != nil {
		return err
	}
	return nil
}
