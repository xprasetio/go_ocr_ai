package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"ocr_ai/config"
)

type D1Client struct {
	endpoint   string
	token      string
	httpClient *http.Client
}

type d1Request struct {
	SQL    string `json:"sql"`
	Params []any  `json:"params,omitempty"`
}

type d1Response struct {
	Success bool `json:"success"`
	Errors  []struct {
		Message string `json:"message"`
	} `json:"errors"`
	Result []struct {
		Results []map[string]any `json:"results"`
	} `json:"result"`
}

func NewD1Client(cfg *config.Config) *D1Client {
	return &D1Client{
		endpoint: fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/d1/database/%s/query", cfg.CloudflareAccountID, cfg.CloudflareD1DatabaseID),
		token:    cfg.CloudflareAPIToken,
		httpClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

func (c *D1Client) Exec(ctx context.Context, sql string, params ...any) ([]map[string]any, error) {
	payload, _ := json.Marshal(d1Request{SQL: sql, Params: params})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var out d1Response
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	if !out.Success {
		if len(out.Errors) > 0 {
			return nil, fmt.Errorf("d1 error: %s", out.Errors[0].Message)
		}
		return nil, fmt.Errorf("d1 query failed")
	}

	if len(out.Result) == 0 {
		return nil, nil
	}
	return out.Result[0].Results, nil
}
