package requests

import (
	"reflect"

	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models"
	"sigs.k8s.io/yaml"
)

type ChartRelease struct {
	Name        string      `json:"name"`
	ReleaseName string      `json:"release_name"`
	Version     string      `json:"version"`
	Values      interface{} `json:"values"`
	Namespace   string      `json:"namespace"`
}

func (c ChartRelease) TransformToModels() (models.ChartRelease, error) {
	values, err := yaml.Marshal(c.Values)
	releaseModels := models.ChartRelease{
		Name:        c.Name,
		ReleaseName: c.ReleaseName,
		Version:     c.Version,
		Values:      string(values),
		Namespace:   c.Namespace,
	}
	return releaseModels, err
}

func (c ChartRelease) IsEmpty() bool {
	return reflect.DeepEqual(c, ChartRelease{})
}
