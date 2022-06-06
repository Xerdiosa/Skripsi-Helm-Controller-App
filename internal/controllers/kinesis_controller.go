package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/helpers"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/services"
	"sigs.k8s.io/yaml"
)

type KinesisController struct {
	kinesisService services.IKinesisService
}

func InitKinesisController(kinesisService services.IKinesisService) KinesisController {
	chartController := KinesisController{}
	chartController.kinesisService = kinesisService
	return chartController
}

func (h *KinesisController) Release(res http.ResponseWriter, req *http.Request) {
	requestBody := models.Kinesis{}
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

	err = h.kinesisService.InstallOrUpgradeKinesis(requestBody)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	helpers.Response(res, 200, nil, "success", "-")
}

func (h *KinesisController) UpdateRelease(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	requestBody := models.Kinesis{}
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

	requestBody.Name = vars["kinesis-name"]

	if requestBody.IsEmpty() {
		helpers.Response(res, 400, nil, "error", "cannot process empty request")
		return
	}

	err = h.kinesisService.InstallOrUpgradeKinesis(requestBody)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	helpers.Response(res, 200, nil, "success", "-")
}

func (h *KinesisController) GetAllReleaseName(res http.ResponseWriter, req *http.Request) {
	result, err := h.kinesisService.GetAllReleaseName()
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	helpers.Response(res, 200, result, "success", "-")
}

func (h *KinesisController) GetReleaseDetail(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	result, err := h.kinesisService.GetReleaseDetail(vars["kinesis-name"])
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	helpers.Response(res, 200, result, "success", "-")
}

func (h *KinesisController) RemoveRelease(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	err := h.kinesisService.RemoveKinesis(vars["kinesis-name"])
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	helpers.Response(res, 200, nil, "success", "-")
}
