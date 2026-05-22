package main

import (
	"context"
	"log"
	"os"
	"strings"

	"ocr_ai/config"
	"ocr_ai/internal/domain/ocr"
	appcontainer "ocr_ai/internal/infrastructure/container"
	httpif "ocr_ai/internal/interfaces/http"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use: "ocr-ai",
		RunE: func(cmd *cobra.Command, args []string) error {

			cwd, _ := os.Getwd()
			log.Printf("CWD: %s", cwd)

			// =========================
			// RAW ENV DEBUG
			// =========================
			log.Println("=== RAW ENV DEBUG ===")
			log.Println("PORT:", os.Getenv("PORT"))
			log.Println("CLOUDFLARE_ACCOUNT_ID:", os.Getenv("CLOUDFLARE_ACCOUNT_ID"))
			log.Println("CLOUDFLARE_D1_DATABASE_ID:", os.Getenv("CLOUDFLARE_D1_DATABASE_ID"))
			log.Println("CLOUDFLARE_API_TOKEN length:", len(os.Getenv("CLOUDFLARE_API_TOKEN")))
			log.Println("======================")

			cfg, err := config.Load()
			if err != nil {
				log.Printf("Config Load Error: %v", err)
				return err
			}

			// =========================
			// CONFIG DEBUG
			// =========================
			log.Println("=== CONFIG DEBUG ===")
			log.Println("Config Port:", cfg.Port)
			log.Println("Config CF Account:", cfg.CloudflareAccountID)
			log.Println("Config CF DB:", cfg.CloudflareD1DatabaseID)
			log.Println("Config CF Token length:", len(cfg.CloudflareAPIToken))
			log.Println("====================")

			// ENV LIST DEBUG
			log.Println("--- Environment Variables ---")
			for _, env := range os.Environ() {
				if strings.HasPrefix(env, "CLOUDFLARE") ||
					strings.HasPrefix(env, "PORT") {
					log.Printf("ENV: %s", env)
				}
			}
			log.Println("-----------------------------")

			container, err := appcontainer.Build(cfg)
			if err != nil {
				return err
			}

			usecase := container.Get("ocr_usecase").(*ocr.Usecase)

			if err := usecase.Migrate(context.Background()); err != nil {
				return err
			}

			e := httpif.NewRouter(cfg, usecase)

			log.Printf("server running on :%s", cfg.Port)

			return e.Start(":" + cfg.Port)
		},
	}

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
