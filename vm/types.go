package vms

import (
	_ "libvirt.org/go/libvirt"
)

type UUID string

type ControlInfra struct {
	Cores      []Core         // 모든 코어를 관리
	AliveVM    []*VMInfo      //현재 가동중인 VM의 정보
	VMLocation map[UUID]*Core //UUID를 기반으로 어떤 VM이 어느 Core에 있는지 확인하는 포인터
	// => 여기에다가는 사실 core에 ip주소를 기반으로 어느 ip에 어느 UUID를 가진 VM이 있는지만 확인할 수 있으면 될거같아요.
	// => 근데 이거 자료형을 어떤걸 써야할지 모르겠어서 일단 이렇게 map[UUID]*Core로 해놓은거지 2차원 벡터로 해도 무방할거 같습니다.
	// => 이렇게 한 이유는 이전 버전에서 VMPool map[UUID]*VM 이렇게 해놓았던데 이렇게 하면 VM이 많아졌을 때 너무 오래 걸릴거 같아서 IP를 기반으로 먼저 찾고 해당 core로 넘어가서 처리하면 더 좋지않을까? 생각했어요.
}

type Core struct {
	IP          string           // Core의 IP 주소
	CoreInfoIdx CoreInfo         //코어의 물리적인 정보
	IsAlive     bool             //Core가 살았는지 죽었는지 확인
	VMInfoIdx   map[UUID]*VMInfo // core 안에 VM 정보
	FreeMemory  int              //할당되지 않은 Core의 Memory 자원
	FreeCPU     int              //할당되지 않은 Core의 CPU 자원
}

type CoreInfo struct {
	Memory int
	Cpu    int
	Disk   int
}

// VM 정보
type VMInfo struct {
	IP_VM  string
	UUID   UUID
	Memory int
	Cpu    int
	Disk   int
}

// ----------------------------------------------------
// // 모든 VM의 상태를 관리
// type InfraContext struct {
// 	Computers       []Computer   //관리중인 컴퓨터의 목록
// 	VMPoolAllocated []*VM        //할당된 VM의 목록
// 	VMPool          map[UUID]*VM //UUID를 활용한 VM의 전체 맵
// }

// //모든 컴퓨터, vm에 대한 정보를 담고 있음.
// // json이나 yaml에 지속적으로 업데이트해서 재부팅시에도 유지하는 것 필수

// type InfraManage interface {
// 	UpdateList()
// }

// type VM struct {
// 	VMInfo      VMInfo
// 	IsAlive     bool
// 	IsAllocated bool
// 	IsLocatedAt Computer
// }

// // VM 내부의 상세 정보 표시
// // 사용자가 원하는 하드웨어 정보를 받으면 이 구조체에 저장 후 사용
// type VMInfo struct {
// 	MaxMem    uint64
// 	Memory    uint64
// 	NrVirtCpu uint
// 	CpuTime   uint64
// 	UUID      UUID
// 	IP        string
// }

// // Core의 정보들이 저장되어있음
// type Computer struct {
// 	Name      string
// 	Allocated []VM
// 	IP        string
// 	MAC       string
// 	IsAlive   bool
// }

//
