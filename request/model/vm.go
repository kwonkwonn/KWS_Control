package model

import (
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	vms "github.com/easy-cloud-Knet/KWS_Control/structure"
)

type HardwareInfo struct {
	CPU    uint32 `json:"cpu"`
	Memory uint32 `json:"memory"` // MiB
	Disk   uint32 `json:"disk"`   // MiB
}

type UserInfoVM struct {
	Name              string   `json:"name"`
	Groups            string   `json:"groups"`
	Password          string   `json:"passWord"`
	SSHAuthorizedKeys []string `json:"ssh"`
}

type CreateVMRequest struct {
	DomType      string         `json:"domType"`
	DomName      string         `json:"domName"`
	UUID         structure.UUID `json:"uuid"`
	OS           string         `json:"os"`
	HardwareInfo HardwareInfo   `json:"HWInfo"`
	NetConf      NetDefine      `json:"network"`
	Users        []UserInfoVM   `json:"users"`
}

type DomainDeleteType uint

const (
	HardDelete DomainDeleteType = iota
	SoftDelete
)

type DeleteVMRequest struct {
	UUID structure.UUID   `json:"UUID"`
	Type DomainDeleteType `json:"DeleteType"`
}

type StatusDataType uint

const (
	CpuInfo StatusDataType = iota
	MemInfo
	DiskInfoHi
	SystemInfoHi
)

type GetMachineStatusRequest struct {
	HostDataType StatusDataType `json:"host_dataType"`
}

type NetDefine struct {
	Ips     []string `json:"ips"`
	NetType NetType  `json:"NetType"`
}

type NetType uint

type CreateVMResponse struct {
	State     string `json:"state"`
	MaxMem    uint64 `json:"maxmem"`
	Memory    uint64 `json:"memory"`
	NrVirtCpu uint   `json:"nrVirtCpu"`
	CpuTime   uint64 `json:"cpuTime"`
}

type DeleteVMResponse struct {
}

type CoreMachineCpuInfoResponse struct {
	System float64 `json:"system_time"`
	Idle   float64 `json:"idle_time"`
	Usage  float64 `json:"usage_percent"`
}

type CoreMachineMemoryInfoResponse struct {
	Total       uint64  `json:"total_gb"`
	Used        uint64  `json:"used_gb"`
	Available   uint64  `json:"available_gb"`
	UsedPercent float64 `json:"used_percent"`
}

type CoreMachineDiskInfoResponse struct {
	Total       uint64  `json:"total_gb"`
	Used        uint64  `json:"used_gb"`
	Free        uint64  `json:"free_gb"`
	UsedPercent float64 `json:"used_percent"`
}

type CoreMachineSystemInfoResponse struct {
	Uptime   uint64  `json:"uptime_seconds"`
	BootTime uint64  `json:"boot_time_epoch"`
	CPUTemp  float64 `json:"cpu_temperature,omitempty"`
	RAMTemp  float64 `json:"ram_temperature,omitempty"`
}

type ForceShutdownVMRequest struct {
	UUID vms.UUID `json:"UUID"`
}

type ForceShutdownVMResponse struct {
	Message string `json:"message"`
}

type GetVMStatusRequest struct {
	UUID     structure.UUID `json:"UUID"`
	DataType StatusDataType `json:"dataType"`
}

type Redis struct {
	UUID   structure.UUID `json:"UUID"`
	Status string         `json:"status"`
}

// type instaceStatus int

// const (
// 	termincate instaceStatus = iota
// 	booting
// )
