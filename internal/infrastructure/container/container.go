package container

import (
	"ocr_ai/config"
	"ocr_ai/internal/domain/ocr"
	"ocr_ai/internal/infrastructure/database"
	"ocr_ai/internal/infrastructure/repository"

	"github.com/sarulabs/di/v2"
)

func Build(cfg *config.Config) (di.Container, error) {
	builder, _ := di.NewBuilder()
	_ = builder.Add(di.Def{Name: "config", Build: func(ctn di.Container) (any, error) { return cfg, nil }})
	_ = builder.Add(di.Def{Name: "d1", Build: func(ctn di.Container) (any, error) { return database.NewD1Client(cfg), nil }})
	_ = builder.Add(di.Def{Name: "ocr_repo", Build: func(ctn di.Container) (any, error) {
		return repository.NewOCRRepository(ctn.Get("d1").(*database.D1Client)), nil
	}})
	_ = builder.Add(di.Def{Name: "gemini", Build: func(ctn di.Container) (any, error) { return repository.NewGeminiAnalyzer(cfg), nil }})
	_ = builder.Add(di.Def{Name: "storage", Build: func(ctn di.Container) (any, error) { return repository.NewSupabaseStorage(cfg), nil }})
	_ = builder.Add(di.Def{Name: "ocr_usecase", Build: func(ctn di.Container) (any, error) {
		return ocr.NewUsecase(ctn.Get("ocr_repo").(ocr.Repository), ctn.Get("gemini").(ocr.Analyzer), ctn.Get("storage").(ocr.Storage)), nil
	}})
	return builder.Build(), nil
}
