package api

import (
	"encoding/json"
	"github.com/easy-cloud-Knet/KWS_Control/api/model"
	"github.com/easy-cloud-Knet/KWS_Control/service"
	"net/http"
)

func (c *handlerContext) createVm(w http.ResponseWriter, r *http.Request) {
	err := service.CreateVM(w, r, c.context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("VM created successfully"))
}

func (c *handlerContext) deleteVm(w http.ResponseWriter, r *http.Request) {
	var req model.ApiDeleteVmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := service.DeleteVM(req.UUID, c.context)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // TODO: 코어가 없는 경우 처리
		return
	}

	w.WriteHeader(http.StatusOK)
}
