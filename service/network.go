package service

import (
	"fmt"

	"github.com/easy-cloud-Knet/KWS_Control/client"
	pkgnetwork "github.com/easy-cloud-Knet/KWS_Control/pkg/network"
	vms "github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

// AddCmsSubnet은 기존 VM의 IP가 속한 서브넷으로 CMS에 새 인스턴스 할당을 요청
func AddCmsSubnet(c *client.CmsClient, ctx *vms.ControlContext, uuid vms.UUID) (*client.CmsNewInstanceResponse, error) {
	log := util.GetLogger()

	ip, err := GetVMIPByUUID(ctx, uuid)
	if err != nil {
		log.Error("AddCmsSubnet : GetVMIPByUUID: %v", err)
		return nil, fmt.Errorf("AddCmsSubnet: failed to get VM IP: %w", err)
	}
	subnet, err := pkgnetwork.GetSubnetFromIP(ip)
	if err != nil {
		log.Error("AddCmsSubnet : GetSubnetFromIP: %v", err)
		return nil, fmt.Errorf("AddCmsSubnet: failed to get subnet: %w", err)
	}
	resp, err := c.RequestNewInstance(subnet)
	if err != nil {
		log.Error("AddCmsSubnet : RequestNewInstance: %v", err)
		return nil, fmt.Errorf("AddCmsSubnet: CMS request failed: %w", err)
	}
	return resp, nil
}

// NewCmsSubnet은 다음 서브넷을 계산하여 DB에 선점한 뒤 CMS에 인스턴스 할당을 요청
func NewCmsSubnet(c *client.CmsClient, ctx *vms.ControlContext) (*client.CmsNewInstanceResponse, error) {
	log := util.GetLogger()

	lastSubnet := ctx.Last_subnet
	nextLastSubnet := pkgnetwork.FindSubnet(lastSubnet)
	log.Info("NewCmsSubnet : next_last_subnet: %s", nextLastSubnet)

	//CMS 호출 전에 다음 서브넷을 선점하여 동시 호출 시 중복 할당 방지
	_, err := ctx.DB.Exec("UPDATE subnet SET last_subnet = ? WHERE id = 1", nextLastSubnet)
	if err != nil {
		log.Error("Failed to update last_subnet in database: %v", err)
		return nil, fmt.Errorf("NewCmsSubnet: failed to update last_subnet in DB: %w", err)
	}
	ctx.Last_subnet = nextLastSubnet

	resp, err := c.RequestNewInstance(nextLastSubnet)
	if err != nil {
		log.Error("NewCmsSubnet : RequestNewInstance: %v", err)
		// CMS 호출 실패 시 DB를 원래 값으로 롤백
		if _, rbErr := ctx.DB.Exec("UPDATE subnet SET last_subnet = ? WHERE id = 1", lastSubnet); rbErr != nil {
			log.Error("NewCmsSubnet : failed to rollback last_subnet: %v", rbErr)
			return nil, fmt.Errorf("NewCmsSubnet: CMS request failed: %w, rollback also failed: %v", err, rbErr)
		}
		ctx.Last_subnet = lastSubnet
		return nil, fmt.Errorf("NewCmsSubnet: CMS request failed: %w", err)
	}
	return resp, nil
}
func DeleteCmsSubnet(c *client.CmsClient, ctx *vms.ControlContext, uuid vms.UUID) error {
	log := util.GetLogger()

	ip, err := GetVMIPByUUID(ctx, uuid)
	if err != nil {
		log.Error("DeleteCmsSubnet : GetVMIPByUUID: %v", err)
		return fmt.Errorf("DeleteCmsSubnet: failed to get VM IP: %w", err)
	}
	_, err = c.RequestDeleteInstance(ip)
	if err != nil {
		log.Error("DeleteCmsSubnet : RequestDeleteInstance: %v", err)
		return fmt.Errorf("DeleteCmsSubnet: CMS request failed: %w", err)
	}
	return nil
}

func GetVMIPByUUID(ctx *vms.ControlContext, uuid vms.UUID) (string, error) {
	ctx.Resources.RLock()
	defer ctx.Resources.RUnlock()

	core, ok := ctx.Resources.VMLocation[uuid]
	if !ok {
		return "", fmt.Errorf("UUID %s not found in VMLocation", uuid)
	}

	vmInfo, ok := core.VMInfoIdx[uuid]
	if !ok {
		return "", fmt.Errorf("VMInfo for UUID %s not found in Core", uuid)
	}

	return vmInfo.IP_VM, nil
}
