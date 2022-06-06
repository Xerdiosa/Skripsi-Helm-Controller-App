package requests

import (
	"reflect"

	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models"
)

type Module struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Values  string `json:"values"`
	Spec    string `json:"spec"`
}

func (m Module) TransformToModels() models.Module {
	module := models.Module{
		Name:    m.Name,
		Version: m.Version,
		Values:  m.Values,
		Spec:    m.Spec,
	}
	return module
}

func (m Module) IsEmpty() bool {
	return reflect.DeepEqual(m, Module{})
}

const (
	DELETE = "delete"
	KEEP   = "keep"
)

type ModuleRelease struct {
	Name         string `json:"release"`
	ModuleName   string `json:"module"`
	Version      string `json:"version"`
	Values       string `json:"values"`
	DeleteOnFail string `json:"on_failure"`
}

func (m ModuleRelease) TransformToModels(deleteOnFailDefault bool) (models.Module, models.ModuleRelease, bool) {
	module := models.Module{
		Name:    m.ModuleName,
		Version: m.Version,
	}
	ModuleRelease := models.ModuleRelease{
		Name:    m.Name,
		Version: m.Version,
		Values:  m.Values,
	}
	var deleteOnFail bool
	switch m.DeleteOnFail {
	case DELETE:
		deleteOnFail = true
	case KEEP:
		deleteOnFail = false
	default:
		deleteOnFail = deleteOnFailDefault
	}
	return module, ModuleRelease, deleteOnFail
}

func (m ModuleRelease) IsEmpty() bool {
	return reflect.DeepEqual(m, ModuleRelease{})
}
