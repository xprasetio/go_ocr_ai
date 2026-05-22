package http

import (
	"ocr_ai/config"
	"ocr_ai/internal/domain/ocr"
	"ocr_ai/internal/interfaces/http/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(cfg *config.Config, usecase *ocr.Usecase) *echo.Echo {
	e := echo.New()

	// Middleware: Logger & Recover agar error tercetak di console
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{cfg.FrontendAllowedOrigin, "http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.OPTIONS},
	}))

	h := handler.NewOCRHandler(usecase)

	e.GET("/health", func(c echo.Context) error { return c.JSON(200, map[string]string{"status": "ok"}) })

	// Grouping rute
	api := e.Group("/api")
	{
		api.POST("/ocr/analyze", h.Analyze)
		api.PUT("/ocr/:id", h.Update)
		api.GET("/ocr/:id", h.Detail)
		api.GET("/ocr", h.List)
		api.GET("/ocr/export/csv", h.ExportCSV)

		// Alias /api/expenses agar sesuai pemanggilan Axios di log user
		api.GET("/expenses", h.List)
	}

	// Alias /expenses langsung di root agar sesuai pemanggilan Axios di log user
	e.GET("/expenses", h.List)

	return e
}
