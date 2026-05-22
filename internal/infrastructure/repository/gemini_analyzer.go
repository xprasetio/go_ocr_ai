package repository

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"

	"ocr_ai/config"
	"ocr_ai/internal/domain/ocr"
)

type GeminiAnalyzer struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewGeminiAnalyzer(cfg *config.Config) *GeminiAnalyzer {
	return &GeminiAnalyzer{cfg: cfg, httpClient: &http.Client{Timeout: 90 * time.Second}}
}

func (g *GeminiAnalyzer) Analyze(ctx context.Context, fileName, mimeType string, content []byte) (*ocr.Expense, error) {
	prompt := `analisa gambar yang diinputkan sebagai dokumen keuangan (resi/struk/invoice). Solusikan juga kasus dokumen buruk: miring, gelap, blur, atau bukan resi. Jika bukan resi/invoice, isi title "Dokumen tidak valid" dan note alasannya. Jika ada angka 0 dibelakang titik kurang dari tiga maka hapus, contoh 10000.0 menjadi 10000. Balas hanya JSON valid dengan struktur: {"amount":0,"change":0,"created_at":"string","date_shopping":"string","id":0,"items":[{"created_at":"string","expenses_id":0,"id":0,"name":"string","price":0,"qty":0,"subtotal":0,"updated_at":"string"}],"note":"string","parsed_data":"string","receipt_image":"string","title":"string","updated_at":"string"}`
	
	payload := map[string]any{
		"contents": []any{
			map[string]any{
				"parts": []any{
					map[string]any{"text": prompt},
					map[string]any{
						"inline_data": map[string]any{
							"mime_type": mimeType,
							"data":      base64.StdEncoding.EncodeToString(content),
						},
					},
				},
			},
		},
		"generationConfig": map[string]any{
			"response_mime_type": "application/json",
		},
	}

	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.cfg.GeminiModel, g.cfg.GeminiAPIKey)

	var respBody []byte
	var lastStatus int
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		res, err := g.httpClient.Do(req)
		if err != nil {
			return nil, err
		}

		respBody, _ = io.ReadAll(res.Body)
		res.Body.Close()
		lastStatus = res.StatusCode

		if lastStatus == 200 {
			break
		}

		if lastStatus == 503 || lastStatus == 429 {
			wait := time.Duration(math.Pow(2, float64(i))) * time.Second
			log.Printf("Gemini busy (%d), retrying in %v... (Attempt %d/%d)", lastStatus, wait, i+1, maxRetries)
			time.Sleep(wait)
			continue
		}

		log.Printf("Gemini API Error Body: %s", string(respBody))
		return nil, fmt.Errorf("gemini api returned status %d", lastStatus)
	}

	if lastStatus != 200 {
		return nil, fmt.Errorf("gemini unavailable after %d retries (last status %d)", maxRetries, lastStatus)
	}

	// Debug: Print raw response for troubleshooting
	log.Printf("Gemini Raw Response (first 500 chars): %s", string(respBody[:min(500, len(respBody))]))

	var out struct {
		Candidates []struct {
			FinishReason string `json:"finishReason"`
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(respBody, &out); err != nil {
		log.Printf("Failed to unmarshal Gemini response: %v", err)
		return nil, err
	}

	if len(out.Candidates) == 0 || len(out.Candidates[0].Content.Parts) == 0 {
		reason := "UNKNOWN"
		if len(out.Candidates) > 0 {
			reason = out.Candidates[0].FinishReason
		}
		log.Printf("Gemini empty response. Raw: %s", string(respBody))
		return nil, fmt.Errorf("gemini returned empty response (FinishReason: %s)", reason)
	}

	text := cleanJSON(out.Candidates[0].Content.Parts[0].Text)
	var expense ocr.Expense
	if err := json.Unmarshal([]byte(text), &expense); err != nil {
		log.Printf("Failed to unmarshal Gemini JSON: %s", text)
		return nil, err
	}
	
	expense.ParsedData = text
	if expense.Title == "" {
		expense.Title = strings.TrimSuffix(fileName, ".pdf")
	}
	return &expense, nil
}

func cleanJSON(input string) string {
	re := regexp.MustCompile("(?s)```(?:json)?\\s*(.*?)\\s*```")
	if m := re.FindStringSubmatch(input); len(m) == 2 {
		return m[1]
	}
	return strings.TrimSpace(input)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
