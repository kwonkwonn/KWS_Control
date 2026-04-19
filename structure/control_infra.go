package structure

import (
	"database/sql"

	"github.com/easy-cloud-Knet/KWS_Control/util"
)

type ControlContext struct {
	VMRepo      VMRepository
	Resources   *ResourceManager
	Config      Config
	DB          *sql.DB // subnet 등 직접 쿼리용 (추후 별도 Repository로 분리 예정)
	GuacDB      *sql.DB
	Last_subnet string
}

func (c *ControlContext) FindCoreByVmUUID(uuid UUID) *Core {
	log := util.GetLogger()

	coreIdx, err := c.VMRepo.GetInstanceLocation(uuid)
	if err != nil {
		log.Error("Core not found for VM UUID %s", uuid, true)
		return nil
	}
	if coreIdx < 0 || coreIdx >= len(c.Resources.Cores) {
		log.Error("Core index %d out of range for VM UUID %s", coreIdx, uuid, true)
		return nil
	}
	return &c.Resources.Cores[coreIdx]
}
