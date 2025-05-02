package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Control/request"
	"github.com/easy-cloud-Knet/KWS_Control/request/model"

	vms "github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
)

// 새 VM 만드는 무언가.
// 자원 많이 남은 코어를 찾고, 리소스 할당 업데이트, ControlContext 상태 업데이트.
func CreateVM(w http.ResponseWriter, r *http.Request, contextStruct *vms.ControlContext) error {
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

	// vm uuid 생성
	newUUID := vms.UUID(uuid.NewString())
	log.Infof("new UUID: %s", newUUID)

	// ip, err := contextStruct.AssignInternalAddress()
	vmIP := "10.0.0.0" // 할당된 ip 받아오도록 하는 거 필요.

	newVM := &vms.VMInfo{
		UUID:   newUUID,
		Memory: req.HardwareInfo.Memory,
		Cpu:    req.HardwareInfo.CPU,
		Disk:   req.HardwareInfo.Disk,
		IP_VM:  vmIP,
	}

	// selected core 상태 업데이트
	if selectedCore.VMInfoIdx == nil {
		selectedCore.VMInfoIdx = make(map[vms.UUID]*vms.VMInfo)
	}
	selectedCore.VMInfoIdx[newUUID] = newVM
	selectedCore.FreeMemory -= req.HardwareInfo.Memory
	selectedCore.FreeCPU -= req.HardwareInfo.CPU
	selectedCore.FreeDisk -= req.HardwareInfo.Disk
	log.Infof("core %s updated: FreeMemory=%d, FreeCPU=%d, FreeDisk=%d", selectedCore.IP, selectedCore.FreeMemory, selectedCore.FreeCPU, selectedCore.FreeDisk)

	// ControlContext global 상태 업데이트
	if contextStruct.VMLocation == nil {
		contextStruct.VMLocation = make(map[vms.UUID]*vms.Core)
	}
	contextStruct.VMLocation[newUUID] = &contextStruct.Cores[selectedCoreIndex]
	contextStruct.AliveVM = append(contextStruct.AliveVM, newVM)
	log.Infof("VM %s added to ControlContext", newUUID)

	// request.go 부분 필요

	log.Infof("UUID %s CreateVM request success on core %s", newUUID, selectedCore.IP)
	return nil
}

func DeleteVM(uuid vms.UUID, contextStruct *vms.ControlContext) error {
	core := contextStruct.FindCoreByVmUUID(uuid)
	if core == nil {
		return errors.New("core not found")
	}

	client := request.NewCoreClient(core)
	_, err := client.DeleteVM(context.Background(), model.DeleteVMRequest{
		UUID: uuid,
		Type: model.HardDelete,
	})

	return err
}
