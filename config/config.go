package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	GeminiAPIKey           string
	GeminiModel            string
	CloudflareAccountID    string
	CloudflareD1DatabaseID string
	CloudflareAPIToken     string
	SupabaseURL            string
	SupabaseServiceRoleKey string
	SupabaseBucket         string
	FrontendAllowedOrigin  string
}

func Load() (*Config, error) {

	// local development only
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	cfg := &Config{
		Port:                   getEnv("PORT", "8080"),
		GeminiAPIKey:           os.Getenv("GEMINI_API_KEY"),
		GeminiModel:            getEnv("GEMINI_MODEL", "gemini-2.5-flash"),
		CloudflareAccountID:    os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		CloudflareD1DatabaseID: os.Getenv("CLOUDFLARE_D1_DATABASE_ID"),
		CloudflareAPIToken:     os.Getenv("CLOUDFLARE_API_TOKEN"),
		SupabaseURL:            os.Getenv("SUPABASE_URL"),
		SupabaseServiceRoleKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
		SupabaseBucket:         getEnv("SUPABASE_BUCKET", "ocr_ai_receipt"),
		FrontendAllowedOrigin: getEnv(
			"FRONTEND_ALLOWED_ORIGIN",
			"https://ocr-ai.xprasetio.workers.dev",
		),
	}

	// debug
	log.Println("Cloudflare Account ID:", cfg.CloudflareAccountID)
	log.Println("Cloudflare D1 Database ID:", cfg.CloudflareD1DatabaseID)
	log.Println("Cloudflare API Token Length:", len(cfg.CloudflareAPIToken))

	return cfg, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
