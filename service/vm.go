package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/easy-cloud-Knet/KWS_Control/request"
	"github.com/easy-cloud-Knet/KWS_Control/request/model"
	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"

	vms "github.com/easy-cloud-Knet/KWS_Control/structure"
)

// 새 VM 만드는 무언가.
// 자원 많이 남은 코어를 찾고, 리소스 할당 업데이트, ControlContext 상태 업데이트.
func CreateVM(w http.ResponseWriter, r *http.Request, contextStruct *vms.ControlContext) error {
	log := util.GetLogger()

	var req model.CreateVMRequest
	defer r.Body.Close() // defer << 에러가 발생해도 body가 닫히도록 보장.

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("err req body parsing: %v", err, true)
		return errors.New("err req body parsing: " + err.Error())
	}

	log.Info("func CreateVM() memory=%d GiB, cpu=%d, disk=%d GiB", req.HardwareInfo.Memory, req.HardwareInfo.CPU, req.HardwareInfo.Disk, true)

	// err = validateCreateVMRequest(req)

	// guacamole 부분 필요

	// 적합한 코어 찾기
	log.DebugInfo("core selection process. req: memory=%d GiB, cpu=%d, disk=%d",
		req.HardwareInfo.Memory, req.HardwareInfo.CPU, req.HardwareInfo.Disk)

	var selectedCore *vms.Core = nil
	selectedCoreIndex := -1
	aliveCount := 0

	for i := range contextStruct.Cores {
		core := &contextStruct.Cores[i]
		log.DebugInfo("core %s checking: FreeMemory=%d, FreeCPU=%d, FreeDisk=%d, IsAlive=%t", core.IP, core.FreeMemory, core.FreeCPU, core.FreeDisk, core.IsAlive)

		if !core.IsAlive {
			log.DebugWarn("core %s is not alive, skipping", core.IP)
			continue
		}
		aliveCount++

		memoryOk := core.FreeMemory >= req.HardwareInfo.Memory
		cpuOk := core.FreeCPU >= req.HardwareInfo.CPU
		diskOk := core.FreeDisk >= req.HardwareInfo.Disk

		if !memoryOk {
			log.DebugWarn("%s !memoryOk: req=%d, available=%d", core.IP, req.HardwareInfo.Memory, core.FreeMemory)
		}
		if !cpuOk {
			log.DebugWarn("%s !cpuOk   : req=%d, available=%d", core.IP, req.HardwareInfo.CPU, core.FreeCPU)
		}
		if !diskOk {
			log.DebugWarn("%s !diskOk  : req=%d, available=%d", core.IP, req.HardwareInfo.Disk, core.FreeDisk)
		}

		if memoryOk && cpuOk && diskOk {
			selectedCore = core
			selectedCoreIndex = i
			log.DebugInfo("core found: %s", selectedCore.IP)
			break
		} else {
			log.DebugInfo("%s rejected: memory=%t, cpu=%t, disk=%t", core.IP, memoryOk, cpuOk, diskOk)
		}
	}

	if selectedCore == nil {
		log.Error("No suitable core found! Total cores: %d, Alive cores: %d, Required: Memory=%d CPU=%d Disk=%d",
			len(contextStruct.Cores), aliveCount, req.HardwareInfo.Memory, req.HardwareInfo.CPU, req.HardwareInfo.Disk, true)

		if aliveCount > 0 {
			log.DebugError("alive cores:")
			for i := range contextStruct.Cores {
				core := &contextStruct.Cores[i]
				if core.IsAlive {
					log.DebugError("  %s: Memory=%d/%d, CPU=%d/%d, Disk=%d/%d",
						core.IP, core.FreeMemory, core.CoreInfoIdx.Memory, core.FreeCPU, core.CoreInfoIdx.Cpu, core.FreeDisk, core.CoreInfoIdx.Disk)
				}
			}
		} else {
			log.DebugError("no alive cores available")
		}

		return errors.New("selectedCore == nil")
	}
	var privateKeyPEM, publicKeyOpenSSH, err = GenerateSshKey()
	if err != nil {
		log.Error("GenerateSshKey() failed: %v", err, true)
		return err
	}

	// 문제가 생겼을 때 지우는 무언가
	var guacamoleConfigured = false
	var coreResourcesAllocated = false
	var uuid = structure.UUID(req.UUID.String().(string))

	cleanup := func() {
		if guacamoleConfigured {
			log.Info("clean up clean up")
			if cleanupErr := CleanupGuacamoleConfig(string(uuid), contextStruct.GuacDB); cleanupErr != nil {
				log.Error("Failed to cleanup Guacamole config during rollback: %v", cleanupErr)
			}
		}
		if coreResourcesAllocated {
			delete(selectedCore.VMInfoIdx, uuid)
			selectedCore.FreeMemory += req.HardwareInfo.Memory
			selectedCore.FreeCPU += req.HardwareInfo.CPU
			selectedCore.FreeDisk += req.HardwareInfo.Disk
		}
	}

	instanceIp, err := contextStruct.AssignInternalAddress()
	if err != nil {
		log.Error("AssignInternalAddress() failed: %v", err, true)
		return err
	}

	fmt.Printf("AssignInternalAddress(): %s", instanceIp)
	fmt.Println(publicKeyOpenSSH) // TODO: 코어로 보내줘야함

	userPass := GuacamoleConfig(req.Users[0].Name, string(uuid), instanceIp, privateKeyPEM, contextStruct.GuacDB)

	if userPass == "" {
		log.Error("Failed to configure Guacamole", true)
		return errors.New("failed to configure Guacamole")
	}
	guacamoleConfigured = true

	newVM := &structure.VMInfo{
		UUID:         uuid,
		GuacPassword: userPass,
		Memory:       req.HardwareInfo.Memory,
		Cpu:          req.HardwareInfo.CPU,
		Disk:         req.HardwareInfo.Disk,
		IP_VM:        instanceIp,
	}

	// selected core 상태 업데이트
	if selectedCore.VMInfoIdx == nil {
		selectedCore.VMInfoIdx = make(map[vms.UUID]*vms.VMInfo)
	}
	selectedCore.VMInfoIdx[uuid] = newVM
	selectedCore.FreeMemory -= req.HardwareInfo.Memory
	selectedCore.FreeCPU -= req.HardwareInfo.CPU
	selectedCore.FreeDisk -= req.HardwareInfo.Disk
	coreResourcesAllocated = true

	log.DebugInfo("core %s updated: FreeMemory=%d, FreeCPU=%d, FreeDisk=%d", selectedCore.IP, selectedCore.FreeMemory, selectedCore.FreeCPU, selectedCore.FreeDisk)

	req.NetConf.Ips = []string{instanceIp}
	req.NetConf.NetType = 0
	req.Users[0].SSHAuthorizedKeys = []string{publicKeyOpenSSH}

	client := request.NewCoreClient(selectedCore)
	_, err = client.CreateVM(context.Background(), req)
	if err != nil {
		log.Error("Error creating VM on core %s: %v", selectedCore.IP, err, true)
		cleanup() // 직접 지우지 말고 요 함수 하나로--
		return err
	}
	err = contextStruct.AddInstance(newVM, selectedCoreIndex)
	if err != nil {
		log.Error("Error database instance insertion failed: %v", err, true)
		cleanup() // 직접 지우지 말고 요 함수 하나로--
	}

	// ControlContext global 상태 업데이트
	if contextStruct.VMLocation == nil {
		contextStruct.VMLocation = make(map[vms.UUID]*vms.Core)
	}
	contextStruct.VMLocation[uuid] = &contextStruct.Cores[selectedCoreIndex]
	contextStruct.AliveVM = append(contextStruct.AliveVM, newVM)
	log.Info("VM %s added to ControlContext", uuid, true)

	log.Info("UUID %s CreateVM request success on core %s", uuid, selectedCore.IP, true)
	return nil
}

func DeleteVM(uuid vms.UUID, contextStruct *vms.ControlContext) error {
	log := util.GetLogger()
	core := contextStruct.FindCoreByVmUUID(uuid)
	if core == nil {
		log.Error("VM with UUID %s not found", string(uuid))
		return fmt.Errorf("VM with UUID %s not found", string(uuid))
	}

	client := request.NewCoreClient(core)
	_, err := client.DeleteVM(context.Background(), model.DeleteVMRequest{
		UUID: uuid,
		Type: model.HardDelete,
	})
	if err != nil {
		log.Error("error deleting VM %s on core %s: %w", uuid, core.IP, err)
		return err
	}

	err = contextStruct.DeleteInstance(uuid)
	if err != nil {
		log.Error("error deleting instance %s from ControlContext: %v", uuid, err)
		return err
	}

	return nil
}

func ShutdownVM(uuid vms.UUID, contextStruct *vms.ControlContext) error {
	core := contextStruct.FindCoreByVmUUID(uuid)
	if core == nil {
		return fmt.Errorf("VM with UUID %s not found", string(uuid))
	}

	client := request.NewCoreClient(core)
	_, err := client.ForceShutdownVM(context.Background(), model.ForceShutdownVMRequest{
		UUID: uuid,
	})

	if err != nil {
		return err
	}

	foundIndex := -1
	for i, vm := range contextStruct.AliveVM {
		if vm.UUID == uuid {
			foundIndex = i
			break
		}
	}

	if foundIndex != -1 {
		contextStruct.AliveVM = slices.Delete(contextStruct.AliveVM, foundIndex, foundIndex+1)
	}

	return nil
}

func GetVMCpuInfo(uuid vms.UUID, contextStruct *vms.ControlContext) (model.CoreMachineCpuInfoResponse, error) {
	log := util.GetLogger()

	core := contextStruct.FindCoreByVmUUID(uuid)
	if core == nil {
		msg := fmt.Sprintf("VM with UUID %s not found", string(uuid))
		log.Error(msg, true)
		return model.CoreMachineCpuInfoResponse{}, errors.New(msg)
	}

	client := request.NewCoreClient(core)

	cpuInfo, err := client.GetVMCpuInfo(context.Background(), uuid)
	if err != nil {
		msg := fmt.Sprintf("Error getting CPU info for VM %s on core %s: %v", uuid, core.IP, err)
		log.Error(msg, true)
		return model.CoreMachineCpuInfoResponse{}, errors.New(msg)
	}

	log.DebugInfo("Retrieved CPU status for VM %s on core %s", uuid, core.IP)
	return cpuInfo, nil
}

func GetVMMemoryInfo(uuid vms.UUID, contextStruct *vms.ControlContext) (model.CoreMachineMemoryInfoResponse, error) {
	log := util.GetLogger()

	core := contextStruct.FindCoreByVmUUID(uuid)
	if core == nil {
		msg := fmt.Sprintf("VM with UUID %s not found", string(uuid))
		log.Error(msg, true)
		return model.CoreMachineMemoryInfoResponse{}, errors.New(msg)
	}

	client := request.NewCoreClient(core)

	memoryInfo, err := client.GetVMMemoryInfo(context.Background(), uuid)
	if err != nil {
		msg := fmt.Sprintf("Error getting memory info for VM %s on core %s: %v", uuid, core.IP, err)
		log.Error(msg, true)
		return model.CoreMachineMemoryInfoResponse{}, errors.New(msg)
	}

	log.DebugInfo("Retrieved Memory status for VM %s on core %s", uuid, core.IP)
	return memoryInfo, nil
}

func GetVMDiskInfo(uuid vms.UUID, contextStruct *vms.ControlContext) (model.CoreMachineDiskInfoResponse, error) {
	log := util.GetLogger()

	core := contextStruct.FindCoreByVmUUID(uuid)
	if core == nil {
		msg := fmt.Sprintf("VM with UUID %s not found", string(uuid))
		log.Error(msg, true)
		return model.CoreMachineDiskInfoResponse{}, errors.New(msg)
	}

	client := request.NewCoreClient(core)

	diskInfo, err := client.GetVMDiskInfo(context.Background(), uuid)
	if err != nil {
		msg := fmt.Sprintf("Error getting disk info for VM %s on core %s: %v", uuid, core.IP, err)
		log.Error(msg, true)
		return model.CoreMachineDiskInfoResponse{}, errors.New(msg)
	}

	log.DebugInfo("Retrieved Disk status for VM %s on core %s", uuid, core.IP)
	return diskInfo, nil
}
