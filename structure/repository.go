package structure

// VM 인스턴스의 데이터베이스 영속성 인터페이스
// ControlContext에서 DB 관련 책임을 분리
type VMRepository interface {
	AddInstance(instanceInfo *VMInfo, coreIdx int) error
	UpdateInstance(instanceInfo *VMInfo) error
	DeleteInstance(uuid UUID) error
	GetInstance(uuid UUID) (*VMInfo, error)
	GetInstanceLocation(uuid UUID) (int, error)
	GetAllInstanceInfo() ([]VMInfo, []int, error)
}
