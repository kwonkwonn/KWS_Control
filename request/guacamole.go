package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/sirupsen/logrus"
)

type GuacamoleClient struct {
	baseURL   string
	client    *http.Client
	authToken string
}

func (c *GuacamoleClient) AuthToken() string {
	return c.authToken
}

func NewGuacamoleClient(config *structure.Config) *GuacamoleClient {
	return &GuacamoleClient{
		baseURL: config.GuacBaseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type AuthenticateResponse struct {
	AuthToken  string   `json:"authToken"`
	Username   string   `json:"username"`
	DataSource string   `json:"dataSource"`
	Available  []string `json:"availableDataSources"`
}

func (c *GuacamoleClient) Authenticate(ctx context.Context, username, password string) error {
	// Authenticate는 w-form-www-urlencoded 형식으로 요청을 보내서 doRequest를 사용하지 않음

	log := logrus.New()
	log.SetReportCaller(true)

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/tokens", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send auth request: %w", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var authResponse AuthenticateResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	c.authToken = authResponse.AuthToken
	log.Info("Successfully authenticated with Guacamole")

	return nil
}

func (c *GuacamoleClient) doRequest(ctx context.Context, method, path string, requestBody, responseBody interface{}) error {
	log := logrus.New()
	log.SetReportCaller(true)

	var reqBodyReader io.Reader
	if requestBody != nil {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		log.Println("(-> Guacamole) request body:", string(jsonData))
		reqBodyReader = bytes.NewBuffer(jsonData)
	}

	requestURL := c.baseURL + path
	if c.authToken != "" {
		if strings.Contains(path, "?") {
			requestURL += "&token=" + c.authToken
		} else {
			requestURL += "?token=" + c.authToken
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, reqBodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	if responseBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	return nil
}
