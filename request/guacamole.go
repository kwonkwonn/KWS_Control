package request

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/easy-cloud-Knet/KWS_Control/structure"
	"github.com/easy-cloud-Knet/KWS_Control/util"
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
	log := util.GetLogger()
	log.Println("NewGuacamoleClient (-> Guacamole) baseURL:", config.GuacBaseURL)
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

	log := util.GetLogger()

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/tokens", strings.NewReader(data.Encode()))
	log.Println("(-> Guacamole) request body:", data.Encode())
	log.Println("(-> Guacamole) request URL:", c.baseURL+"/api/tokens")
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send auth request to %s: %w", c.baseURL, err)
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
	log.Info("Successfully authenticated with Guacamole", true)

	return nil
}
