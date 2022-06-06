package responses

type ChartRelease struct {
	Name        string      `json:"name"`
	ReleaseName string      `json:"release_name"`
	Version     string      `json:"version"`
	Values      interface{} `json:"values"`
	Revision    int         `json:"revision"`
	Namespace   string      `json:"namespace"`
}
