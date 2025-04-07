package structure

type UUID string

type Config struct {
	VmInternalSubnets []string `yaml:"vm_internal_subnets"`
	Cores             []string `yaml:"cores"`
	Port              int      `yaml:"port"`
}

type ControlInfra struct {
	Config     Config
	Cores      []Core         // 모든 코어를 관리
	AliveVM    []*VMInfo      //현재 가동중인 VM의 정보
	VMLocation map[UUID]*Core //UUID를 기반으로 어떤 VM이 어느 Core에 있는지 확인하는 포인터
	// => 여기에다가는 사실 core에 ip주소를 기반으로 어느 ip에 어느 UUID를 가진 VM이 있는지만 확인할 수 있으면 될거같아요.
	// => 근데 이거 자료형을 어떤걸 써야할지 모르겠어서 일단 이렇게 map[UUID]*Core로 해놓은거지 2차원 벡터로 해도 무방할거 같습니다.
	// => 이렇게 한 이유는 이전 버전에서 VMPool map[UUID]*VM 이렇게 해놓았던데 이렇게 하면 VM이 많아졌을 때 너무 오래 걸릴거 같아서 IP를 기반으로 먼저 찾고 해당 core로 넘어가서 처리하면 더 좋지않을까? 생각했어요.
}

// memory: GiB
// disk: GiB
// CPU: Logical Core num

type Core struct {
	IP          string           // 코어의 IP 주소
	Port        uint16           // 코어의 포트 번호
	CoreInfoIdx CoreInfo         //코어의 물리적인 정보
	IsAlive     bool             //Core가 살았는지 죽었는지 확인
	VMInfoIdx   map[UUID]*VMInfo // core 안에 VM 정보
	FreeMemory  uint16           //할당되지 않은 코어의 Memory 자원 GiB
	FreeCPU     uint8            //할당되지 않은 코어의 CPU 자원 논리 코어 수
	FreeDisk    uint16           //할당되지 않은 코어의 Disk 자원 GiB
}

type CoreInfo struct {
	Memory uint16 // 코어의 전체 메모리 GiB
	Cpu    uint8  // 코어의 전체 CPU 논리 코어 수
	Disk   uint16 // 코어의 전체 디스크 GiB
}

type VMInfo struct {
	IP_VM  string
	UUID   UUID
	Memory uint16 // VM의 메모리 GiB
	Cpu    uint8  // VM의 CPU 논리 코어 수
	Disk   uint16 // VM의 디스크 GiB
}
