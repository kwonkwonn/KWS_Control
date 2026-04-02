package api

import (
	"encoding/json"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Control/service"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

type ApiDeleteVmRequest struct {
	UUID structure.UUID `json:"uuid"`
}

func (c *handlerContext) deleteVm(w http.ResponseWriter, r *http.Request) {
	log := util.GetLogger()
	defer r.Body.Close()

	var req ApiDeleteVmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("deleteVm: failed to decode request body: %v", err, true)
		util.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := service.DeleteVM(req.UUID, c.context, c.rdb)
	if err != nil {
		log.Error("deleteVm: failed to delete VM: %v", err, true)
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
