package api

import (
	"encoding/json"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Control/api/model"
	"github.com/easy-cloud-Knet/KWS_Control/service"
	"github.com/sirupsen/logrus"
)

func (c *handlerContext) createVm(w http.ResponseWriter, r *http.Request) {
	log := logrus.New()
	log.SetReportCaller(true)

	err := service.CreateVM(w, r, c.context)
	if err != nil {
		h := w.Header()
		h.Del("Content-Length")
		h.Set("Content-Type", "text/plain; charset=utf-8")
		h.Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusMethodNotAllowed)

		log.Errorf("Failed to create VM: %v", err)

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
