package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gotrics-node/internal/hardware"
	"gotrics-node/internal/metrics"
	"net/http"
	"time"
)

type Sender struct {
	client    *http.Client
	serverURL string
	authToken string
}

func NewSender(url string, token string) *Sender {
	return &Sender{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		serverURL: url,
		authToken: token,
	}
}

func (s *Sender) Send(m *metrics.Metrics) error {
	payload, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	req, err := http.NewRequest("POST", s.serverURL+"/api/v1/metrics", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.authToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("server returned non-success status: %s", resp.Status)
	}

	return nil
}

func (s *Sender) SendHardwareInfo(info *hardware.HardwareInfo) error {
	payload, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal hardware info: %w", err)
	}

	req, err := http.NewRequest("POST", s.serverURL+"/api/v1/hardware", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create hardware request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.authToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send hardware request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("server returned non-success status for hardware: %s", resp.Status)
	}

	return nil
}
