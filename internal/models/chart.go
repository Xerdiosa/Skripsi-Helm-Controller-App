package models

import (
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models/responses"
)

type ChartRelease struct {
	Model
	ModuleReleaseID uint   `json:"-"`
	Name            string `json:"name"`
	ReleaseName     string `json:"release_name"`
	Version         string `json:"version"`
	Values          string `json:"values"`
	Revision        int    `json:"revision"`
	Namespace       string `json:"namespace"`
}

func (c ChartRelease) TransformToResponse() responses.ChartRelease {
	response := responses.ChartRelease{
		Name:        c.Name,
		ReleaseName: c.ReleaseName,
		Version:     c.Version,
		Values:      c.Values,
		Revision:    c.Revision,
		Namespace:   c.Namespace,
	}
	return response
}
