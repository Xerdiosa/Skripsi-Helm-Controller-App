package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/helpers"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models/requests"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/services"
	"sigs.k8s.io/yaml"
)

type ChartController struct {
	chartService services.IChartService
}

func InitChartController(helmService services.IChartService) ChartController {
	chartController := ChartController{}
	chartController.chartService = helmService
	return chartController
}

func (h *ChartController) Release(res http.ResponseWriter, req *http.Request) {
	requestBody := requests.ChartRelease{}
	val, err := ioutil.ReadAll(req.Body)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}

	err = yaml.Unmarshal([]byte(val), &requestBody)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	if requestBody.IsEmpty() {
		helpers.Response(res, 400, nil, "error", "cannot process empty request")
		return
	}

	request, err := requestBody.TransformToModels()
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}

	err = h.chartService.InstallOrUpgradeChart(request)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	helpers.Response(res, 200, nil, "success", "-")
}

func (h *ChartController) UpdateRelease(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	requestBody := requests.ChartRelease{}
	val, err := ioutil.ReadAll(req.Body)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}

	err = yaml.Unmarshal([]byte(val), &requestBody)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	if requestBody.IsEmpty() {
		helpers.Response(res, 400, nil, "error", "cannot process empty request")
		return
	}

	request, err := requestBody.TransformToModels()
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}

	request.ReleaseName = vars["chart-name"]

	err = h.chartService.InstallOrUpgradeChart(request)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	helpers.Response(res, 200, nil, "success", "-")
}

func (h *ChartController) GetAllReleaseName(res http.ResponseWriter, req *http.Request) {
	result, err := h.chartService.GetAllReleaseName()
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	helpers.Response(res, 200, result, "success", "-")
}

func (h *ChartController) GetReleaseDetail(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	result, err := h.chartService.GetReleaseDetail(vars["chart-name"])
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	helpers.Response(res, 200, result.TransformToResponse(), "success", "-")
}

func (h *ChartController) RemoveRelease(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	err := h.chartService.RemoveChart(vars["chart-name"])
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	helpers.Response(res, 200, nil, "success", "-")
}
