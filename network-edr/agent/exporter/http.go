package exporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/example/network-edr/agent/flow"
)

type HTTP struct {
	BaseURL string
	Client  *http.Client
}

func NewHTTP(base string) *HTTP {
	return &HTTP{
		BaseURL: strings.TrimRight(base, "/"),
		Client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (h *HTTP) Send(ctx context.Context, f flow.Flow) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.BaseURL+"/api/v1/ingest", bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("server returned %s", resp.Status)
	}
	return nil
}
