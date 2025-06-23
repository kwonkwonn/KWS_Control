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
	defer r.Body.Close()

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("VM created successfully"))
}

func (c *handlerContext) deleteVm(w http.ResponseWriter, r *http.Request) {
	var req model.ApiDeleteVmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := service.DeleteVM(req.UUID, c.context)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // TODO: 코어가 없는 경우 처리
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *handlerContext) shutdownVm(w http.ResponseWriter, r *http.Request) {
	var req model.ApiShutdownVmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err := service.ShutdownVM(req.UUID, c.context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *handlerContext) vmStatus(w http.ResponseWriter, r *http.Request) {
	log := logrus.New()
	log.SetReportCaller(true)

	var req model.ApiVmStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Errorf("Failed to decode request body: %v", err)
		return
	}
	defer r.Body.Close()

	statusType := req.Type
	if statusType != "cpu" && statusType != "memory" && statusType != "disk" {
		http.Error(w, "Invalid status type. Must be 'cpu', 'memory', or 'disk'", http.StatusBadRequest)
		return
	}

	var data any
	var err error

	switch statusType {
	case "cpu":
		data, err = service.GetVMCpuInfo(req.UUID, c.context)
	case "memory":
		data, err = service.GetVMMemoryInfo(req.UUID, c.context)
	case "disk":
		data, err = service.GetVMDiskInfo(req.UUID, c.context)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf("Failed to get VM status: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (c *handlerContext) vmConnect(w http.ResponseWriter, r *http.Request) {
	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	log := logrus.New()
	log.SetReportCaller(true)

	var req model.ApiVmConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Errorf("Failed to decode request body: %v", err)
		return
	}

	authToken, err := service.GetGuacamoleToken(req.UUID, c.context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf("Failed to get Guacamole token: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"authToken": authToken}); err != nil {
		log.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
