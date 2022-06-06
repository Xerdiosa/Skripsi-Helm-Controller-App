package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models/responses"
)

func Response(w http.ResponseWriter, httpStatus int, data interface{}, status string, errorMsg string) {
	apiResponse := responses.APIResponse{}
	apiResponse.Data = data
	apiResponse.Status = status
	apiResponse.Error = errorMsg

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(apiResponse)
}
