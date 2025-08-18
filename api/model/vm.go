package model

import "github.com/easy-cloud-Knet/KWS_Control/structure"

type ApiDeleteVmRequest struct {
	UUID structure.UUID `json:"uuid"`
}

type ApiShutdownVmRequest struct {
	UUID structure.UUID `json:"uuid"`
}

type ApiVmStatusRequest struct {
	UUID structure.UUID `json:"uuid"`
	Type string         `json:"type"` // "cpu", "memory", or "disk"
}

type ApiVmConnectRequest struct {
	UUID structure.UUID `json:"uuid"`
}

// 요건 core쪽이랑 이야기 맞춰야--
// request/model/vm.go 와 동일하게 유지
const (
	VMStatusBooting    = "booting"
	VMStatusRunning    = "running"
	VMStatusStopped    = "stopped"
	VMStatusTerminated = "terminated"
	VMStatusUnknown    = "unknown"
)

type Redis struct {
	UUID   structure.UUID `json:"UUID"`
	Status string         `json:"status"`
}

func ValidateAndNormalizeStatus(status string) string {
	if status == "" || status == "null" {
		return VMStatusUnknown
	}

	switch status {
	case VMStatusBooting, VMStatusRunning, VMStatusStopped, VMStatusTerminated, VMStatusUnknown:
		return status
	default:
		return VMStatusUnknown
	}
}
