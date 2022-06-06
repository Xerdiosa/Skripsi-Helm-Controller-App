package responses

type APIResponse struct {
	Data   interface{} `json:"data"`
	Error  string      `json:"error"`
	Status string      `json:"status"`
}
