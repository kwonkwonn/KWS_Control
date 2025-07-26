package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

type subnet_parse struct {
	val1 string
	val2 string
	val3 string
}

func Find_subnet(last_subnet string) string {

	value := make([]int, 3)

	for i := 0; i < 3; i++ {
		var j int
		var temp string
		for last_subnet[j] != '.' {
			temp = temp + string(last_subnet[j])
			j++
		}
		value[i], _ = strconv.Atoi(temp)
	}

	if value[2] >= 255 {
		if value[1] >= 255 {
			if value[0] >= 255 {
				return ""
			} else {
				value[0]++
				value[1] = 0
				value[2] = 0
			}
		} else {
			value[1]++
			value[2] = 0
		}
	} else {
		value[2]++
	}

	result := fmt.Sprintf("%s.%s.%s", strconv.Itoa(value[0]), strconv.Itoa(value[1]), strconv.Itoa(value[2]))
	return result
}
