package api

import (
	"encoding/json"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Control/service"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

type ApiShutdownVmRequest struct {
	UUID structure.UUID `json:"uuid"`
}

type ApiShutdownVmResponse struct {
	Message string `json:"message"`
}

type ApiForceShutdownVmRequest struct {
	UUID structure.UUID `json:"uuid"`
}

type ApiForceShutdownVmResponse struct {
	Message string `json:"message"`
}

func (c *handlerContext) shutdownVm(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req ApiShutdownVmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := service.ShutdownVM(req.UUID, c.context, c.rdb)
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
