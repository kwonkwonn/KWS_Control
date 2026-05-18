package structure

import (
	"slices"
	"sync"
)

// ResourceManager는 코어/VM의 런타임 인메모리 상태를 관리
// ControlContext에서 뮤텍스와 상태 필드를 분리
type ResourceManager struct {
	mu         sync.RWMutex
	Cores      []Core         // 모든 코어를 리스트
	AliveVM    []*VMInfo      // 현재 가동중인 VM의 정보
	VMLocation map[UUID]*Core // UUID 기반 VM  코어 위치 확인
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		VMLocation: make(map[UUID]*Core),
	}
}

func (rm *ResourceManager) Lock()    { rm.mu.Lock() }
func (rm *ResourceManager) Unlock()  { rm.mu.Unlock() }
func (rm *ResourceManager) RLock()   { rm.mu.RLock() }
func (rm *ResourceManager) RUnlock() { rm.mu.RUnlock() }

// HardwareRequirement는 코어 선택 시 필요한 자원 요구량 구조체
type HardwareRequirement struct {
	Memory uint32 // MiB
	CPU    uint32 // logical cores
	Disk   uint32 // MiB
}

// CoreSelectionResult는 SelectCore의 반환값으로 진단 정보를 내포
type CoreSelectionResult struct {
	Core       *Core
	Index      int
	AliveCount int
	TotalCores int
}

// SelectCore는 요청한 자원을 만족하는 첫 번째 살아있는 코어를 탐색.
// RLock 범위 내에서 코어 슬라이스를 순회 (네트워크 콜은 호출자가 락 밖에서 수행).
// 적합한 코어가 없으면 Core==nil로 반환하며, 진단 로그를 위한 카운트 정보를 함께 제공
func (rm *ResourceManager) SelectCore(req HardwareRequirement) CoreSelectionResult {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	result := CoreSelectionResult{
		Index:      -1,
		TotalCores: len(rm.Cores),
	}

	for i := range rm.Cores {
		core := &rm.Cores[i]
		if !core.IsAlive {
			continue
		}
		result.AliveCount++

		if core.FreeMemory >= req.Memory && core.FreeCPU >= req.CPU && core.FreeDisk >= req.Disk {
			result.Core = core
			result.Index = i
			return result
		}
	}
	return result
}

// AllocateResources는 코어의 VMInfoIdx 맵에 VM을 등록하고 Free* 필드를 차감
func (rm *ResourceManager) AllocateResources(core *Core, uuid UUID, vm *VMInfo, req HardwareRequirement) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if core.VMInfoIdx == nil {
		core.VMInfoIdx = make(map[UUID]*VMInfo)
	}
	core.VMInfoIdx[uuid] = vm
	core.FreeMemory -= req.Memory
	core.FreeCPU -= req.CPU
	core.FreeDisk -= req.Disk
}

// DeallocateResources는 AllocateResources의 역연산
func (rm *ResourceManager) DeallocateResources(core *Core, uuid UUID, req HardwareRequirement) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	delete(core.VMInfoIdx, uuid)
	core.FreeMemory += req.Memory
	core.FreeCPU += req.CPU
	core.FreeDisk += req.Disk
}

// RegisterVM은 VMLocation 맵과 AliveVM 슬라이스에 VM을 동시에 등록
func (rm *ResourceManager) RegisterVM(uuid UUID, core *Core, vm *VMInfo) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	rm.VMLocation[uuid] = core
	rm.AliveVM = append(rm.AliveVM, vm)
}

// UnregisterAlive는 AliveVM 슬라이스에서 해당 UUID를 제거
func (rm *ResourceManager) UnregisterAlive(uuid UUID) bool {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	for i, vm := range rm.AliveVM {
		if vm.UUID == uuid {
			rm.AliveVM = slices.Delete(rm.AliveVM, i, i+1)
			return true
		}
	}
	return false
}
