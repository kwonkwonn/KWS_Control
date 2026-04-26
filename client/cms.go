package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/easy-cloud-Knet/KWS_Control/util"
)

// CmsClient는 CMS서비스에 서브넷/인스턴스 생성을 요청하는 HTTP 클라이언트.
type CmsClient struct {
	baseURL string
	client  *http.Client
}

type CmsResponse struct {
	IP      string `json:"ip"`
	MacAddr string `json:"macAddr"`
	SdnUUID string `json:"sdnUUID"`
}

type cmsRequestBody struct {
	Subnet string `json:"Subnet"`
}

func NewCmsClient() *CmsClient {
	host := os.Getenv("CMS_HOST")
	if host == "" {
		log := util.GetLogger()
		log.Error("CMS_HOST Re:Check your env variable", true)
		host = "localhost:8080"
		log.Warn("CMS_HOST set: %s", host, true)
	}
	return &CmsClient{
		baseURL: host,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RequestNewInstance는 주어진 서브넷에 대해 CMS에 새 인스턴스 할당을 요청한다.
func (c *CmsClient) RequestNewInstance(subnet string) (*CmsResponse, error) {
	log := util.GetLogger()

	reqURL := fmt.Sprintf("http://%s/New/Instance", c.baseURL)
	jsonBody, err := json.Marshal(cmsRequestBody{Subnet: subnet})
	if err != nil {
		log.Error("CMS : failed to marshal JSON: %v", err)
		return nil, fmt.Errorf("CmsClient.RequestNewInstance: failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error("CMS : failed to NewRequest: %v", err)
		return nil, fmt.Errorf("CmsClient.RequestNewInstance: failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	log.DebugInfo("Making request to: %s", reqURL)
	log.DebugInfo("Request body: %s", string(jsonBody))

	resp, err := c.client.Do(req)
	if err != nil {
		log.Error("CMS : failed to send request: %v", err)
		return nil, fmt.Errorf("CmsClient.RequestNewInstance: failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("CMS : CMS returned status: %s", resp.Status)
		return nil, fmt.Errorf("CMS server returned non-OK status: %s", resp.Status)
	}

	var addrResp CmsResponse
	if err := json.NewDecoder(resp.Body).Decode(&addrResp); err != nil {
		log.Error("CMS : failed to decode CMS response: %v", err)
		return nil, fmt.Errorf("CmsClient.RequestNewInstance: failed to decode response: %w", err)
	}
	return &addrResp, nil
}
