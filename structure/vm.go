package structure

type UUID string

type Config struct {
	VmInternalSubnets []string `yaml:"vm_internal_subnets"`
	Cores             []string `yaml:"cores"`
	Port              int      `yaml:"port"`
	DB                string   `yaml:"db"`
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
