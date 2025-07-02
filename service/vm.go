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

	vms "github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/sirupsen/logrus"
)

// 새 VM 만드는 무언가.
// 자원 많이 남은 코어를 찾고, 리소스 할당 업데이트, ControlContext 상태 업데이트.
func CreateVM(w http.ResponseWriter, r *http.Request, contextStruct *vms.ControlContext) error {
	//var err error
	log := logrus.New()
	log.SetReportCaller(true)

	var req model.CreateVMRequest
	defer r.Body.Close() // defer << 에러가 발생해도 body가 닫히도록 보장.

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("err req body parsing: %v", err)
		return errors.New("err req body parsing: " + err.Error())
	}

	log.Infof("func CreateVM() memory=%dMiB, cpu=%d, disk=%dMiB", req.HardwareInfo.Memory, req.HardwareInfo.CPU, req.HardwareInfo.Disk)

	// err = validateCreateVMRequest(req)

	// guacamole 부분 필요

	// 적합한 코어 찾기
	var selectedCore *vms.Core = nil
	selectedCoreIndex := -1
	for i := range contextStruct.Cores {
		core := &contextStruct.Cores[i]
		log.Infof("core %s checking: FreeMemory=%d, FreeCPU=%d, FreeDisk=%d, IsAlive=%t", core.IP, core.FreeMemory, core.FreeCPU, core.FreeDisk, core.IsAlive)

		if core.IsAlive && core.FreeMemory >= req.HardwareInfo.Memory && core.FreeCPU >= req.HardwareInfo.CPU && core.FreeDisk >= req.HardwareInfo.Disk {
			selectedCore = core
			selectedCoreIndex = i
			log.Infof("core found: %s", selectedCore.IP)
			break
		}
	}

	if selectedCore == nil {
		log.Errorf("selectedCore == nil")
		return errors.New("selectedCore == nil")
	}
	var privateKeyPEM, publicKeyOpenSSH, err = GenerateSshKey()
	if err != nil {
		log.Errorf("GenerateSshKey() failed: %v", err)
		return err
	}

	instanceIp, err := contextStruct.AssignInternalAddress()
	if err != nil {
		log.Errorf("AssignInternalAddress() failed: %v", err)
		return err
	}
	fmt.Printf("AssignInternalAddress(): %s", instanceIp)
	GuacamoleConfig(req.Users[0].Name, string(req.UUID), instanceIp, privateKeyPEM, contextStruct.Config)
	fmt.Println(publicKeyOpenSSH) // TODO: 코어로 보내줘야함

	newVM := &vms.VMInfo{
		UUID:   req.UUID,
		Memory: req.HardwareInfo.Memory,
		Cpu:    req.HardwareInfo.CPU,
		Disk:   req.HardwareInfo.Disk,
		IP_VM:  instanceIp,
	}

	// selected core 상태 업데이트
	if selectedCore.VMInfoIdx == nil {
		selectedCore.VMInfoIdx = make(map[vms.UUID]*vms.VMInfo)
	}
	selectedCore.VMInfoIdx[req.UUID] = newVM
	selectedCore.FreeMemory -= req.HardwareInfo.Memory
	selectedCore.FreeCPU -= req.HardwareInfo.CPU
	selectedCore.FreeDisk -= req.HardwareInfo.Disk
	log.Infof("core %s updated: FreeMemory=%d, FreeCPU=%d, FreeDisk=%d", selectedCore.IP, selectedCore.FreeMemory, selectedCore.FreeCPU, selectedCore.FreeDisk)

	req.NetConf.Ips = []string{instanceIp}
	req.NetConf.NetType = 0
	req.Users[0].SSHAuthorizedKeys = []string{publicKeyOpenSSH}

	client := request.NewCoreClient(selectedCore)
	_, err = client.CreateVM(context.Background(), req)
	if err != nil {
		delete(selectedCore.VMInfoIdx, req.UUID)
		selectedCore.FreeMemory += req.HardwareInfo.Memory
		selectedCore.FreeCPU += req.HardwareInfo.CPU
		selectedCore.FreeDisk += req.HardwareInfo.Disk

		log.Infof("Error creating VM on core %s: %v", selectedCore.IP, err)
		return err
	}

	// ControlContext global 상태 업데이트
	if contextStruct.VMLocation == nil {
		contextStruct.VMLocation = make(map[vms.UUID]*vms.Core)
	}
	contextStruct.VMLocation[req.UUID] = &contextStruct.Cores[selectedCoreIndex]
	contextStruct.AliveVM = append(contextStruct.AliveVM, newVM)
	log.Infof("VM %s added to ControlContext", req.UUID)

	log.Infof("UUID %s CreateVM request success on core %s", req.UUID, selectedCore.IP)
	return nil
}

func DeleteVM(uuid vms.UUID, contextStruct *vms.ControlContext) error {
	core := contextStruct.FindCoreByVmUUID(uuid)
	if core == nil {
		return fmt.Errorf("VM with UUID %s not found", string(uuid))
	}

	client := request.NewCoreClient(core)
	_, err := client.DeleteVM(context.Background(), model.DeleteVMRequest{
		UUID: uuid,
		Type: model.HardDelete,
	})

	return err
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
	log := logrus.New()
	log.SetReportCaller(true)

	core := contextStruct.FindCoreByVmUUID(uuid)
	if core == nil {
		msg := fmt.Sprintf("VM with UUID %s not found", string(uuid))
		log.Error(msg)
		return model.CoreMachineCpuInfoResponse{}, errors.New(msg)
	}

	client := request.NewCoreClient(core)

	cpuInfo, err := client.GetVMCpuInfo(context.Background(), uuid)
	if err != nil {
		msg := fmt.Sprintf("Error getting CPU info for VM %s on core %s: %v", uuid, core.IP, err)
		log.Error(msg)
		return model.CoreMachineCpuInfoResponse{}, errors.New(msg)
	}

	log.Infof("Retrieved CPU status for VM %s on core %s", uuid, core.IP)
	return cpuInfo, nil
}

func GetVMMemoryInfo(uuid vms.UUID, contextStruct *vms.ControlContext) (model.CoreMachineMemoryInfoResponse, error) {
	log := logrus.New()
	log.SetReportCaller(true)

	core := contextStruct.FindCoreByVmUUID(uuid)
	if core == nil {
		msg := fmt.Sprintf("VM with UUID %s not found", string(uuid))
		log.Error(msg)
		return model.CoreMachineMemoryInfoResponse{}, errors.New(msg)
	}

	client := request.NewCoreClient(core)

	memoryInfo, err := client.GetVMMemoryInfo(context.Background(), uuid)
	if err != nil {
		msg := fmt.Sprintf("Error getting memory info for VM %s on core %s: %v", uuid, core.IP, err)
		log.Error(msg)
		return model.CoreMachineMemoryInfoResponse{}, errors.New(msg)
	}

	log.Infof("Retrieved Memory status for VM %s on core %s", uuid, core.IP)
	return memoryInfo, nil
}

func GetVMDiskInfo(uuid vms.UUID, contextStruct *vms.ControlContext) (model.CoreMachineDiskInfoResponse, error) {
	log := logrus.New()
	log.SetReportCaller(true)

	core := contextStruct.FindCoreByVmUUID(uuid)
	if core == nil {
		msg := fmt.Sprintf("VM with UUID %s not found", string(uuid))
		log.Error(msg)
		return model.CoreMachineDiskInfoResponse{}, errors.New(msg)
	}

	client := request.NewCoreClient(core)

	diskInfo, err := client.GetVMDiskInfo(context.Background(), uuid)
	if err != nil {
		msg := fmt.Sprintf("Error getting disk info for VM %s on core %s: %v", uuid, core.IP, err)
		log.Error(msg)
		return model.CoreMachineDiskInfoResponse{}, errors.New(msg)
	}

	log.Infof("Retrieved Disk status for VM %s on core %s", uuid, core.IP)
	return diskInfo, nil
}
