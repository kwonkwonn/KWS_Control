package service

import (
	vms "github.com/easy-cloud-Knet/KWS_Control/structure"
)

// Service Layer DTO.
// DTOs for service layer to decouple API layer from internal implementation details.
// API Handler will convert API layer DTOs to service layer DTOs.
// Service layer will convert to internal client/model structures as needed.
//

type CreateVMInput struct {
	UUID         vms.UUID
	DomType      string
	DomName      string
	OS           string
	HardwareInfo HardwareSpec
	Network      NetworkSpec
	Users        []UserSpec
	SubnetType   string // "Add" 또는 그 외(=신규 서브넷)
}

type HardwareSpec struct {
	CPU    uint32
	Memory uint32 // MiB
	Disk   uint32 // MiB
}

type NetworkSpec struct {
	IPs []string
}

type UserSpec struct {
	Name              string
	Groups            string
	Password          string
	SSHAuthorizedKeys []string
}

// Output DTOs
type VMInfo struct {
	UUID   vms.UUID
	CPU    uint32
	Memory uint32 // MiB
	Disk   uint32 // MiB
	IP     string
	Status string
}

type VMCpuStatus struct {
	System float64
	Idle   float64
	Usage  float64
}

type VMMemoryStatus struct {
	Total       uint64
	Used        uint64
	Available   uint64
	UsedPercent float64
}

type VMDiskStatus struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64
}
