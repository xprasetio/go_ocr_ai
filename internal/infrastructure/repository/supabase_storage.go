package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"net/http"
	"strings"
	"time"

	"ocr_ai/config"
)

type SupabaseStorage struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewSupabaseStorage(cfg *config.Config) *SupabaseStorage {
	return &SupabaseStorage{cfg: cfg, httpClient: &http.Client{Timeout: 45 * time.Second}}
}

func normalizeSupabaseBaseURL(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", err
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid SUPABASE_URL: %q", raw)
	}
	// Users often paste the REST base like https://<ref>.supabase.co/rest/v1/
	// but Storage endpoints are rooted at the project base URL.
	u.Path = ""
	u.RawQuery = ""
	u.Fragment = ""
	return strings.TrimRight(u.String(), "/"), nil
}

func looksLikeJWT(token string) bool {
	parts := strings.Split(token, ".")
	return len(parts) == 3 && parts[0] != "" && parts[1] != "" && parts[2] != ""
}

func supabaseKeyType(token string) string {
	switch {
	case token == "":
		return "empty"
	case strings.HasPrefix(token, "sb_secret_"):
		return "sb_secret"
	case looksLikeJWT(token):
		return "jwt"
	default:
		return "unknown"
	}
}

func (s *SupabaseStorage) Upload(ctx context.Context, fileName, mimeType string, content []byte) (string, error) {
	baseURL, err := normalizeSupabaseBaseURL(s.cfg.SupabaseURL)
	if err != nil {
		return "", err
	}
	apiKey := strings.TrimSpace(s.cfg.SupabaseServiceRoleKey)
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", baseURL, s.cfg.SupabaseBucket, fileName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(content))
	if err != nil {
		return "", err
	}
	if looksLikeJWT(apiKey) {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Content-Type", mimeType)
	req.Header.Set("x-upsert", "true")
	res, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("supabase upload failed: %s - %s (bucket=%s, key_type=%s)", res.Status, strings.TrimSpace(string(body)), s.cfg.SupabaseBucket, supabaseKeyType(apiKey))
	}
	return fmt.Sprintf("%s/storage/v1/object/public/%s/%s", baseURL, s.cfg.SupabaseBucket, fileName), nil
}
