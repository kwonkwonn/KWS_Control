package WorkerCont

type CoreMachineCpuInfoResp struct {
	System float64 `json:"system_time"`
	Idle   float64 `json:"idle_time"`
	Usage  float64 `json:"usage_percent"`
}

type CoreMachineMemoryInfoResp struct {
	Total       uint64  `json:"total_gb"`
	Used        uint64  `json:"used_gb"`
	Available   uint64  `json:"available_gb"`
	UsedPercent float64 `json:"used_percent"`
}

type CoreMachineDiskInfoResp struct {
	Total       uint64  `json:"total_gb"`
	Used        uint64  `json:"used_gb"`
	Free        uint64  `json:"free_gb"`
	UsedPercent float64 `json:"used_percent"`
}

type CoreMachineSystemInfoResp struct {
	Uptime   uint64  `json:"uptime_seconds"`
	BootTime uint64  `json:"boot_time_epoch"`
	CPUTemp  float64 `json:"cpu_temperature,omitempty"`
	RAMTemp  float64 `json:"ram_temperature,omitempty"`
}
