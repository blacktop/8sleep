package eightsleep

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

const (
	clientAPIURL = "https://client-api.8slp.net/v1"
	appAPIURL    = "https://app-api.8slp.net"
	authURL      = "https://auth-api.8slp.net/v1/tokens"

	knownClientID     = "0894c7f33bb94800a03f1f4df13a4f38"
	knownClientSecret = "f0954a3ed5763ba3d06834c73731a32f15f168f47d4f164751275def86db0c76"

	tokenRefreshBufferSec = 120
	defaultTimeoutSec     = 240

	MIN_TEMP_F = 55
	MAX_TEMP_F = 110
	MIN_TEMP_C = 13
	MAX_TEMP_C = 44
)

var POSSIBLE_SLEEP_STAGES = []string{"bedTimeLevel", "initialSleepLevel", "finalSleepLevel"}

type Client struct {
	mu sync.RWMutex

	email, password string
	tz              *time.Location

	clientID, clientSecret string

	http  *http.Client
	token *Token

	isPod   bool
	hasBase bool

	me      *Profile
	devices []Device
}

func NewClient(email, password, tz string) (*Client, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone %s: %w", tz, err)
	}
	return &Client{
		email:        email,
		password:     password,
		tz:           loc,
		clientID:     knownClientID,
		clientSecret: knownClientSecret,
		http: &http.Client{
			Timeout: time.Second * defaultTimeoutSec,
		},
	}, nil
}

/* -------------------- Public high-level API -------------------- */

func (c *Client) Start(ctx context.Context) error {
	if err := c.refreshToken(ctx); err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	if err := c.fetchProfile(ctx); err != nil {
		return fmt.Errorf("failed to fetch profile: %w", err)
	}
	if err := c.fetchDevices(ctx); err != nil {
		return fmt.Errorf("failed to fetch devices: %w", err)
	}
	return nil
}

func (c *Client) Stop() { /* nothing to close right now */ }

func (c *Client) RoomTemperature(ctx context.Context) (float64, error) {
	panic("not implemented")
	// TODO: get trends and calculate room temperature average (average both sides if both are active)
}

func (c *Client) TurnOn(ctx context.Context) error {
	url := fmt.Sprintf("%s/v1/users/%s/temperature/pod?ignoreDeviceErrors=false", appAPIURL, c.me.ID)
	body := map[string]any{
		"currentState": map[string]string{"type": "smart"},
	}
	var resp TemperatureState
	if err := c.doJSON(ctx, http.MethodPut, url, body, &resp); err != nil {
		return fmt.Errorf("failed to turn on device: %w", err)
	}

	for _, device := range resp.Devices {
		if device.CurrentState.Type == "off" {
			return fmt.Errorf("failed to turn on device %s: %s", device.Device.DeviceID, device.CurrentState.Type)
		}
	}

	return nil
}

func (c *Client) TurnOff(ctx context.Context) error {
	url := fmt.Sprintf("%s/v1/users/%s/temperature/pod?ignoreDeviceErrors=false", appAPIURL, c.me.ID)
	body := map[string]any{
		"currentState": map[string]string{"type": "off"},
	}
	var resp TemperatureState
	if err := c.doJSON(ctx, http.MethodPut, url, body, &resp); err != nil {
		return fmt.Errorf("failed to turn off device: %w", err)
	}

	for _, device := range resp.Devices {
		if device.CurrentState.Type != "off" {
			return fmt.Errorf("failed to turn off device %s: %s", device.Device.DeviceID, device.CurrentState.Type)
		}
	}

	return nil
}

func (c *Client) SetTemperature(ctx context.Context, degrees int, unit UnitOfTemperature) error {
	url := fmt.Sprintf("%s/v1/users/%s/temperature/pod?ignoreDeviceErrors=false", appAPIURL, c.me.ID)
	body := map[string]any{
		"currentLevel": tempToHeatingLevel(degrees, unit),
	}
	var resp TemperatureState
	if err := c.doJSON(ctx, http.MethodPut, url, body, &resp); err != nil {
		return fmt.Errorf("failed to set temperature: %w", err)
	}

	for _, device := range resp.Devices {
		if device.CurrentLevel != tempToHeatingLevel(degrees, unit) {
			return fmt.Errorf("failed to set temperature on device %s: %s", device.Device.DeviceID, device.CurrentState.Type)
		}
	}

	return nil
}

/* -------------------- internal helpers -------------------- */

func (c *Client) Headers() http.Header {
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	h.Set("Accept", "application/json")
	h.Set("Accept-Encoding", "gzip")
	h.Set("User-Agent", "okhttp/4.9.3")
	c.mu.RLock()
	if c.token != nil {
		h.Set("Authorization", "Bearer "+c.token.Bearer)
	}
	c.mu.RUnlock()
	return h
}

func (c *Client) refreshToken(ctx context.Context) error {
	c.mu.RLock()
	needsRefresh := c.token == nil || time.Until(c.token.Expiration) < time.Second*tokenRefreshBufferSec
	c.mu.RUnlock()
	if !needsRefresh {
		return nil
	}

	body := map[string]string{
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"grant_type":    "password",
		"username":      c.email,
		"password":      c.password,
	}
	var res struct {
		AccessToken string  `json:"access_token"`
		ExpiresIn   float64 `json:"expires_in"`
		UserID      string  `json:"userId"`
	}
	if err := c.doJSON(ctx, http.MethodPost, authURL, body, &res); err != nil {
		return err
	}

	c.mu.Lock()
	c.token = &Token{
		Bearer:     res.AccessToken,
		Expiration: time.Now().Add(time.Duration(res.ExpiresIn) * time.Second),
		MainID:     res.UserID,
	}
	c.mu.Unlock()
	return nil
}

func (c *Client) fetchProfile(ctx context.Context) error {
	url := clientAPIURL + "/users/me"
	var data struct {
		User Profile `json:"user"`
	}
	if err := c.doJSON(ctx, http.MethodGet, url, nil, &data); err != nil {
		return err
	}
	c.mu.Lock()
	for _, f := range data.User.Features {
		if f == "cooling" {
			c.isPod = true
		}
		if f == "elevation" {
			c.hasBase = true
		}
	}
	c.me = &data.User
	c.mu.Unlock()
	return nil
}

func (c *Client) fetchDevices(ctx context.Context) error {
	for _, device := range c.me.Devices {
		url := clientAPIURL + "/devices/" + device
		var data struct {
			Result Device `json:"result"`
		}
		if err := c.doJSON(ctx, http.MethodGet, url, nil, &data); err != nil {
			return err
		}
		c.mu.Lock()
		c.devices = append(c.devices, data.Result)
		c.mu.Unlock()
	}
	return nil
}

func (c *Client) doJSON(ctx context.Context, method, url string, payload any, out any) error {
	var body *bytes.Reader

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewReader(b)
	} else {
		body = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = c.Headers()

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute %s request: %w", method, err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if res.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		data, err = io.ReadAll(gzipReader)
		if err != nil {
			return fmt.Errorf("failed to read gzipped response body: %w", err)
		}
	}

	log.Debugf("HTTP %s %s: %d\n%s", method, url, res.StatusCode, string(data))

	return json.NewDecoder(bytes.NewReader(data)).Decode(out)
}
