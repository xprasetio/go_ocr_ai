package main

import (
	"context"
	"log"
	"os"
	"strings"

	"ocr_ai/config"
	appcontainer "ocr_ai/internal/infrastructure/container"
	httpif "ocr_ai/internal/interfaces/http"
	"ocr_ai/internal/domain/ocr"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{Use: "ocr-ai", RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		log.Printf("CWD: %s", cwd)

		cfg, err := config.Load()
		if err != nil {
			log.Printf("Config Load Error: %v", err)
			return err
		}

		log.Println("--- Environment Debug ---")
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, "CLOUDFLARE") {
				log.Printf("Found Env: %s", env)
			}
		}
		log.Printf("CF Token from Config struct: [%s]", cfg.CloudflareAPIToken)
		log.Println("-------------------------")

		container, err := appcontainer.Build(cfg); if err != nil { return err }
		usecase := container.Get("ocr_usecase").(*ocr.Usecase)
		if err := usecase.Migrate(context.Background()); err != nil { return err }
		e := httpif.NewRouter(cfg, usecase)
		log.Printf("server running on :%s", cfg.Port)
		return e.Start(":" + cfg.Port)
	}}
	if err := root.Execute(); err != nil { log.Fatal(err) }
}
