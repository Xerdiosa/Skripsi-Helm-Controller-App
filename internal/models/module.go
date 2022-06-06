package models

type Module struct {
	Model
	Name    string `gorm:"uniqueIndex:module_search" json:"name"`
	Version string `gorm:"uniqueIndex:module_search" json:"version"`
	Values  string `json:"values"`
	Spec    string `json:"spec"`
}

type ModuleRelease struct {
	Model
	ModuleID   uint
	ModuleName string
	Name       string
	Version    string
	Values     string
	Revision   int
}

type ModuleTemplate struct {
	Module  string
	Release string
	Version string
	Values  map[string]interface{}
}

type Spec struct {
	Chart map[string][]interface{} `json:"spec"`
}
