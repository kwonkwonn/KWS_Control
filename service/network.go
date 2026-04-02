package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	pkgnetwork "github.com/easy-cloud-Knet/KWS_Control/pkg/network"
	vms "github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
)

type CmsClient struct {
	baseURL string
	client  *http.Client
}

type CmsResponse struct {
	IP      string `json:"ip"`
	MacAddr string `json:"macAddr"`
	SdnUUID string `json:"sdnUUID"`
}

type CmsRequest struct {
	Subnet string `json:"Subnet"`
}

// fmt.Sprintf("%s/New/Instance", CMS_HOST)
func NewCmsClient() *CmsClient {
	CMS_HOST := os.Getenv("CMS_HOST")
	if CMS_HOST == "" {
		log := util.GetLogger()
		log.Error("CMS_HOST Re:Check your env variable", true)
		CMS_HOST = "localhost:8080"
		log.Warn("CMS_HOST set: %s", CMS_HOST, true)
	}
	return &CmsClient{
		baseURL: CMS_HOST,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *CmsClient) CmsRequest(Subnet string) (*CmsResponse, error) {
	log := util.GetLogger()

	req_url := fmt.Sprintf("http://%s/New/Instance", c.baseURL)
	reqBody := CmsRequest{Subnet: Subnet}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Error("CMS : failed to marshal JSON: %v", err)
		return nil, fmt.Errorf("CmsRequest: failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", req_url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error("CMS : failed to NewRequest: %v", err)
		return nil, fmt.Errorf("CmsRequest: failed to create HTTP request: %w", err)
	}

	// Content-Type 헤더 설정
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	log.DebugInfo("Making request to: %s", req_url)
	log.DebugInfo("Request body: %s", string(jsonBody))

	resp, err := c.client.Do(req)
	if err != nil {
		log.Error("CMS : failed to send request: %v", err)
		return nil, fmt.Errorf("CmsRequest: failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("CMS : CMS returned status: %s", resp.Status)
		return nil, fmt.Errorf("CMS server returned non-OK status: %s", resp.Status)
	}
	var addrResp CmsResponse
	if err := json.NewDecoder(resp.Body).Decode(&addrResp); err != nil {
		log.Error("CMS : failed to decode CMS response: %v", err)
		return nil, fmt.Errorf("CmsRequest: failed to decode response: %w", err)
	}

	return &addrResp, nil
}

func (c *CmsClient) AddCmsSubnet(ctx *vms.ControlContext, uuid vms.UUID) (*CmsResponse, error) {
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
	temp, err := c.CmsRequest(subnet)
	if err != nil {
		log.Error("AddCmsSubnet : c.CmsRequest(subnet): %v", err)
		return nil, fmt.Errorf("AddCmsSubnet: CmsRequest failed: %w", err)
	}

	return temp, nil

}

func (c *CmsClient) NewCmsSubnet(ctx *vms.ControlContext) (*CmsResponse, error) {
	log := util.GetLogger()

	last_subnet := ctx.Last_subnet
	next_last_subnet := pkgnetwork.FindSubnet(last_subnet)
	log.Info("NewCmsSubnet : next_last_subnet: %s", next_last_subnet)

	// DB를 먼저 업데이트하여 서브넷을 선점한다.
	// CMS 호출 전에 선점해야 실패 시 동일 서브넷이 중복 할당되는 것을 방지할 수 있다.
	_, err := ctx.DB.Exec("UPDATE subnet SET last_subnet = ? WHERE id = 1", next_last_subnet)
	if err != nil {
		log.Error("Failed to update last_subnet in database: %v", err)
		return nil, fmt.Errorf("NewCmsSubnet: failed to update last_subnet in DB: %w", err)
	}
	ctx.Last_subnet = next_last_subnet

	temp, err := c.CmsRequest(next_last_subnet)
	if err != nil {
		log.Error("NewCmsSubnet : c.CmsRequest(subnet): %v", err)
		// CMS 호출 실패 시 DB를 원래 값으로 롤백
		if _, rbErr := ctx.DB.Exec("UPDATE subnet SET last_subnet = ? WHERE id = 1", last_subnet); rbErr != nil {
			log.Error("NewCmsSubnet : failed to rollback last_subnet: %v", rbErr)
			return nil, fmt.Errorf("NewCmsSubnet: CmsRequest failed: %w, rollback also failed: %v", err, rbErr)
		}
		ctx.Last_subnet = last_subnet
		return nil, fmt.Errorf("NewCmsSubnet: CmsRequest failed: %w", err)
	}

	return temp, nil
}

func GetVMIPByUUID(ctx *vms.ControlContext, uuid vms.UUID) (string, error) {
	core, ok := ctx.VMLocation[uuid]
	if !ok {
		return "", fmt.Errorf("UUID %s not found in VMLocation", uuid)
	}

	vmInfo, ok := core.VMInfoIdx[uuid]
	if !ok {
		return "", fmt.Errorf("VMInfo for UUID %s not found in Core", uuid)
	}

	return vmInfo.IP_VM, nil
}

