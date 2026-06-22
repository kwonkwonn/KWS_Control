package api

import (
	"encoding/json"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Control/service"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

type ApiVmStatusRequest struct {
	UUID structure.UUID `json:"uuid"`
	Type string         `json:"type"` // "cpu", "memory", or "disk"
}

type ApiVmCpuStatusResponse struct {
	System float64 `json:"system_time"`
	Idle   float64 `json:"idle_time"`
	Usage  float64 `json:"usage_percent"`
}

type ApiVmMemoryStatusResponse struct {
	Total       uint64  `json:"total_gb"`
	Used        uint64  `json:"used_gb"`
	Available   uint64  `json:"available_gb"`
	UsedPercent float64 `json:"used_percent"`
}

type ApiVmDiskStatusResponse struct {
	Total       uint64  `json:"total_gb"`
	Used        uint64  `json:"used_gb"`
	Free        uint64  `json:"free_gb"`
	UsedPercent float64 `json:"used_percent"`
}

func (c *handlerContext) vmStatus(w http.ResponseWriter, r *http.Request) {
	log := util.GetLogger()
	defer r.Body.Close()

	var req ApiVmStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "Invalid request body")
		log.Error("Failed to decode request body: %v", err, true)
		return
	}

	statusType := req.Type
	if statusType != "cpu" && statusType != "memory" && statusType != "disk" {
		util.RespondError(w, http.StatusBadRequest, "Invalid status type. Must be 'cpu', 'memory', or 'disk'")
		return
	}

	var data any
	var err error

	switch statusType {
	case "cpu":
		cpu, e := service.GetVMCpuInfo(req.UUID, c.context)
		data, err = ApiVmCpuStatusResponse{System: cpu.System, Idle: cpu.Idle, Usage: cpu.Usage}, e
	case "memory":
		mem, e := service.GetVMMemoryInfo(req.UUID, c.context)
		data, err = ApiVmMemoryStatusResponse{Total: mem.Total, Used: mem.Used, Available: mem.Available, UsedPercent: mem.UsedPercent}, e
	case "disk":
		disk, e := service.GetVMDiskInfo(req.UUID, c.context)
		data, err = ApiVmDiskStatusResponse{Total: disk.Total, Used: disk.Used, Free: disk.Free, UsedPercent: disk.UsedPercent}, e
	}

	if err != nil {
		log.Error("Failed to get VM status: %v", err, true)
		util.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, data)
}
