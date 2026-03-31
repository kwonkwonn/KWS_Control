package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/easy-cloud-Knet/KWS_Control/client/model"
	"github.com/easy-cloud-Knet/KWS_Control/service"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

func (c *handlerContext) createVm(w http.ResponseWriter, r *http.Request) {
	log := util.GetLogger()
	defer r.Body.Close()

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		log.Warn("No Content-Type header specified, assuming application/json", true)
	} else if !strings.Contains(contentType, "application/json") {
		util.RespondError(w, http.StatusBadRequest, "Content-Type must be application/json")
		return
	}

	var req model.CreateVMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("err req body parsing: %v", err, true)
		if strings.Contains(err.Error(), "invalid character") {
			util.RespondError(w, http.StatusBadRequest, "invalid JSON format in request body")
		} else {
			util.RespondError(w, http.StatusBadRequest, "err req body parsing: "+err.Error())
		}
		return
	}

	if req.HardwareInfo.Memory == 0 || req.HardwareInfo.CPU == 0 || req.HardwareInfo.Disk == 0 {
		util.RespondError(w, http.StatusBadRequest, "Memory, CPU, and Disk must be non-zero")
		return
	}

	err := service.CreateVM(req, c.context, c.rdb)
	if err != nil {
		log.Error("Failed to create VM: %v", err, true)
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("VM created successfully"))
}
