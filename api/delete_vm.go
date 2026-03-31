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
	defer r.Body.Close()

	var req ApiDeleteVmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := service.DeleteVM(req.UUID, c.context, c.rdb)
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
