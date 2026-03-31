package api

import (
	"encoding/json"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Control/service"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

type ApiVmInfoRequest struct {
	UUID structure.UUID `json:"uuid"`
}

type ApiVmInfoResponse struct {
	UUID   structure.UUID `json:"uuid"`
	CPU    uint32         `json:"cpu"`
	Memory uint32         `json:"memory"` // MiB
	Disk   uint32         `json:"disk"`   // MiB
	IP     string         `json:"ip"`
}

func (c *handlerContext) vmInfo(w http.ResponseWriter, r *http.Request) {
	log := util.GetLogger()
	defer r.Body.Close()

	var req ApiVmInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "invalid request body")
		log.Error("Invalid request body: %v", err, true)
		return
	}

	vmInfo, err := service.GetVMInfoFromRedis(r.Context(), c.rdb, req.UUID)
	if err != nil {
		log.Error("failed to get vm info from redis: %v", err, true)
		util.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	response := ApiVmInfoResponse{
		UUID:   vmInfo.UUID,
		CPU:    vmInfo.CPU,
		Memory: vmInfo.Memory,
		Disk:   vmInfo.Disk,
		IP:     vmInfo.IP,
	}

	util.RespondJSON(w, http.StatusOK, response)
	log.Info("retrieved vm info from redis: UUID=%s", string(req.UUID), true)
}
