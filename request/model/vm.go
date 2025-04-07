package model

type HardwareInfo struct {
	CPU    int `json:"cpu"`
	Memory int `json:"memory"`
}

type UserInfoVM struct {
	Name              string   `json:"name"`
	Groups            string   `json:"groups"`
	Password          string   `json:"passWord"`
	SSHAuthorizedKeys []string `json:"ssh"`
}

type CreateVMRequest struct {
	DomType      string       `json:"domType"`
	DomName      string       `json:"domName"`
	UUID         string       `json:"uuid"`
	OS           string       `json:"os"`
	HardwareInfo HardwareInfo `json:"HWInfo"`
	NetConf      NetDefine    `json:"network"`
	Users        []UserInfoVM `json:"users"`
}

type HostDataType uint

const (
	CpuInfo HostDataType = iota
	MemInfo
	DiskInfoHi
	SystemInfoHi
)

type GetMachineStatusRequest struct {
	HostDataType HostDataType `json:"host_dataType"`
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
