package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/helpers"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/models/requests"
	"github.com/gudangada/data-warehouse/warehouse-controller/internal/services"
)

type ModuleController struct {
	moduleService services.IModuleService
}

func InitModuleController(moduleService services.IModuleService) ModuleController {
	moduleController := ModuleController{}
	moduleController.moduleService = moduleService
	return moduleController
}

func (h *ModuleController) AddModule(res http.ResponseWriter, req *http.Request) {
	requestBody := requests.Module{}
	val, err := ioutil.ReadAll(req.Body)
	requestBody.Name = req.Header.Get("NAME")
	requestBody.Version = req.Header.Get("VERSION")
	requestBody.Spec = string(val)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	if requestBody.IsEmpty() {
		helpers.Response(res, 400, nil, "error", "cannot process empty request")
		return
	}
	request := requestBody.TransformToModels()

	err = h.moduleService.InstallModule(request)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}

	helpers.Response(res, 200, nil, "success", "-")
}

func (h *ModuleController) AddModuleRelease(res http.ResponseWriter, req *http.Request) {
	requestBody := requests.ModuleRelease{}
	val, err := ioutil.ReadAll(req.Body)
	requestBody.Name = req.Header.Get("NAME")
	requestBody.ModuleName = req.Header.Get("MODULE_NAME")
	requestBody.Version = req.Header.Get("VERSION")
	requestBody.Values = string(val)
	requestBody.DeleteOnFail = req.Header.Get("ON_FAILURE")
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	if requestBody.IsEmpty() {
		helpers.Response(res, 400, nil, "error", "cannot process empty request")
		return
	}

	module, moduleRelease, deleteOnFail := requestBody.TransformToModels(true)

	err = h.moduleService.ReleaseModule(module, moduleRelease, deleteOnFail)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}

	helpers.Response(res, 200, nil, "success", "-")
}

func (h *ModuleController) UpdateModuleRelease(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	requestBody := requests.ModuleRelease{}
	val, err := ioutil.ReadAll(req.Body)
	requestBody.Name = req.Header.Get("NAME")
	requestBody.ModuleName = req.Header.Get("MODULE_NAME")
	requestBody.Version = req.Header.Get("VERSION")
	requestBody.Values = string(val)
	requestBody.DeleteOnFail = req.Header.Get("ON_FAILURE")
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	if requestBody.IsEmpty() {
		helpers.Response(res, 400, nil, "error", "cannot process empty request")
		return
	}

	module, moduleRelease, deleteOnFail := requestBody.TransformToModels(false)

	if moduleRelease.Name != vars["release-name"] {
		helpers.Response(res, 400, nil, "error", "invalid data")
		return
	}

	err = h.moduleService.UpdateModuleRelease(module, moduleRelease, deleteOnFail)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}

	helpers.Response(res, 200, nil, "success", "-")
}

func (h *ModuleController) DeleteModuleRelease(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	release := models.ModuleRelease{
		Name: vars["release-name"],
	}

	err := h.moduleService.DeleteModuleRelease(release)
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}

	helpers.Response(res, 200, nil, "success", "-")
}

func (h *ModuleController) GetAllReleaseName(res http.ResponseWriter, req *http.Request) {
	result, err := h.moduleService.GetAllReleaseName()
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	helpers.Response(res, 200, result, "success", "-")
}

func (h *ModuleController) GetReleaseDetail(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	result, err := h.moduleService.GetReleaseDetail(vars["release-name"])
	if err != nil {
		helpers.Response(res, 400, nil, "error", err.Error())
		return
	}
	helpers.Response(res, 200, result, "success", "-")
}
