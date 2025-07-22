package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

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
	return &CmsClient{
		baseURL: CMS_HOST,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *CmsClient) NewCmsSubnet(Subnet string) *CmsResponse {
	log := util.GetLogger()
	c.baseURL = fmt.Sprintf("http://%s/New/Instance", c.baseURL)
	reqBody := CmsRequest{Subnet: Subnet}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Error("CMS : failed to marshal JSON: %w", err)
		return nil
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Error("CMS : failed to NewRequest: %w", err)
		return nil
	}
	resp, err := c.client.Do(req)
	if err != nil {
		log.Error("CMS : failed to create request: %w", err)
		return nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("CMS : CMS returned status: %s", resp.Status)
		return nil
	}
	var addrResp CmsResponse
	if err := json.NewDecoder(resp.Body).Decode(&addrResp); err != nil {
		log.Error("CMS : failed to decode CMS response: %w", err)
		return nil
	}

	return &addrResp
}
