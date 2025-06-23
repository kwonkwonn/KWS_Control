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
