package structure

import "sync"

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
