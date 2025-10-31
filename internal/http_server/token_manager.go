// Package http_server
package http_server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/half-nothing/simple-fsd/internal/interfaces/log"
)

type TokenManager struct {
	logger                  log.LoggerInterface
	flushTokenFlushCallback func(flushToken string)
	flushToken              string
	accessToken             string
	flushTimer              *time.Timer
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

func NewTokenManager(
	logger log.LoggerInterface,
	flushToken string,
	flushTokenFlushCallback func(flushToken string),
) *TokenManager {
	manager := &TokenManager{
		logger:                  log.NewLoggerAdapter(logger, "TokenManager"),
		flushToken:              flushToken,
		flushTokenFlushCallback: flushTokenFlushCallback,
		accessToken:             "",
	}
	manager.refreshAccessToken()
	return manager
}

func (m *TokenManager) refreshAccessToken() {
	client := &http.Client{}
	payload := &url.Values{}
	payload.Add("grant_type", "refresh_token")
	payload.Add("refresh_token", m.flushToken)
	payload.Add("scope", "userinfo openid offline_access amdb charts email navdata userdata fmsdata tiles simbrief")
	payload.Add("client_id", "charts-rn-desktop")
	payload.Add("client_secret", "igljsnfBunGqI706JnQRIkQuJB65iscC")
	req, _ := http.NewRequest("POST", "https://identity.api.navigraph.com/connect/token", strings.NewReader(payload.Encode()))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) NavigraphCharts/8.38.3 Chrome/106.0.5249.199 Electron/21.4.1 Safari/537.36")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Host", "identity.api.navigraph.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "zh-CN")
	res, err := client.Do(req)
	if err != nil {
		m.logger.ErrorF("TokenManager refreshAccessToken Error: %s", err.Error())
		return
	}
	if res.StatusCode != http.StatusOK {
		m.logger.ErrorF("TokenManager refreshAccessToken StatusCode: %d", res.StatusCode)
		return
	}
	token := &TokenResponse{}
	data, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	_ = json.Unmarshal(data, token)
	m.accessToken = token.AccessToken
	if m.flushToken != token.RefreshToken {
		m.logger.InfoF("TokenManager refreshAccessToken RefreshToken: %s", m.flushToken)
		m.flushTokenFlushCallback(m.flushToken)
	}
	m.flushTimer = time.AfterFunc(time.Duration(token.ExpiresIn-10)*time.Second, func() { m.refreshAccessToken() })
}

func (m *TokenManager) Shutdown() {
	if m.flushTimer != nil {
		m.flushTimer.Stop()
	}
}

func (m *TokenManager) AccessToken() string {
	return "Bearer " + m.accessToken
}
