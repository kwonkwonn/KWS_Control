package structure

import (
	"github.com/sirupsen/logrus"
)

type ControlContext struct {
	Config     Config
	Cores      []Core         // 모든 코어를 관리
	AliveVM    []*VMInfo      //현재 가동중인 VM의 정보
	VMLocation map[UUID]*Core //UUID를 기반으로 어떤 VM이 어느 Core에 있는지 확인하는 포인터
	// => 여기에다가는 사실 core에 ip주소를 기반으로 어느 ip에 어느 UUID를 가진 VM이 있는지만 확인할 수 있으면 될거같아요.
	// => 근데 이거 자료형을 어떤걸 써야할지 모르겠어서 일단 이렇게 map[UUID]*Core로 해놓은거지 2차원 벡터로 해도 무방할거 같습니다.
	// => 이렇게 한 이유는 이전 버전에서 VMPool map[UUID]*VM 이렇게 해놓았던데 이렇게 하면 VM이 많아졌을 때 너무 오래 걸릴거 같아서 IP를 기반으로 먼저 찾고 해당 core로 넘어가서 처리하면 더 좋지않을까? 생각했어요.
}

func (c *ControlContext) FindCoreByVmUUID(uuid UUID) *Core {

	log := logrus.New()
	log.SetReportCaller(true)

	if core, ok := c.VMLocation[uuid]; ok {
		log.Infof("Found core for VM UUID %s: %s", uuid, core.IP)
		return core
	}
	log.Errorf("Core not found for VM UUID %s", uuid)
	return nil
}
