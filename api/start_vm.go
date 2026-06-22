package api

import (
	"encoding/json"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Control/service"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

type ApiStartVmRequest struct {
	UUID structure.UUID `json:"uuid"`
}

func (c *handlerContext) startVm(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req ApiStartVmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := service.StartVM(req.UUID, c.context)
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
