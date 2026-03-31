package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/easy-cloud-Knet/KWS_Control/service"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

const (
	VMStatusPrepareBegin = "prepare begin"
	VMStatusStartBegin   = "start begin"
	VMStatusStarted      = "started begin"
	VMStatusStopped      = "stopped end"
	VMStatusRelease      = "release end"
	VMStatusMigrate      = "migrate begin"
	VMStatusRestore      = "restort begin"
	VMStatusUnknown      = "unknown"
)

type RedisStatusRequest struct {
	UUID   structure.UUID `json:"UUID"`
	Status string         `json:"status"`
}

func ValidateAndNormalizeStatus(status string) string {
	if status == "" || status == "null" {
		return VMStatusUnknown
	}

	switch status {
	case VMStatusPrepareBegin, VMStatusStartBegin, VMStatusStarted, VMStatusStopped,
		VMStatusRelease, VMStatusMigrate, VMStatusRestore, VMStatusUnknown:
		return status
	default:
		return VMStatusUnknown
	}
}

func (c *handlerContext) redis(w http.ResponseWriter, r *http.Request) {
	log := util.GetLogger()
	defer r.Body.Close()

	var req RedisStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondError(w, http.StatusBadRequest, "Invalid request body")
		log.Error("Invalid request body: %v", err, true)
		return
	}

	originalStatus := req.Status
	normalizedStatus := ValidateAndNormalizeStatus(req.Status)
	if originalStatus != normalizedStatus {
		log.Warn("VM status normalized: UUID=%s, original='%s', normalized='%s'",
			req.UUID, originalStatus, normalizedStatus, true)
	}

	ctx := r.Context()
	err := service.UpdateVMStatusInRedis(ctx, c.rdb, req.UUID, normalizedStatus, time.Now().Unix())
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, "Failed to update VM status in Redis")
		log.Error("Failed to update VM status in Redis: %v", err, true)
		return
	}

	storedValue, err := c.rdb.Get(ctx, string(req.UUID)).Result()
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, "Failed to get value from Redis")
		log.Error("Redis GET failed: %v", err, true)
		return
	}

	log.DebugInfo("Redis 확인 완료 - key: %s, value: %s", req.UUID, storedValue)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("VM status updated in Redis"))
}
