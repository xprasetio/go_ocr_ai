package handler

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"

	"ocr_ai/internal/domain/ocr"

	"github.com/labstack/echo/v4"
)

type OCRHandler struct{ usecase *ocr.Usecase }

func NewOCRHandler(usecase *ocr.Usecase) *OCRHandler { return &OCRHandler{usecase: usecase} }

func (h *OCRHandler) Analyze(c echo.Context) error {
	file, err := c.FormFile("file_attachment")
	if err != nil {
		log.Printf("Error Analyze FormFile: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]any{"message": "file_attachment is required"})
	}
	src, err := file.Open()
	if err != nil {
		log.Printf("Error Analyze Open: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}
	defer src.Close()

	buf := make([]byte, file.Size)
	_, err = src.Read(buf)
	if err != nil {
		log.Printf("Error Analyze Read: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}

	result, err := h.usecase.Analyze(c.Request().Context(), file.Filename, file.Header.Get("Content-Type"), buf)
	if err != nil {
		log.Printf("Error Analyze Usecase: %v", err)
		return c.JSON(http.StatusBadGateway, map[string]any{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (h *OCRHandler) Update(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var payload ocr.Expense
	if err := c.Bind(&payload); err != nil {
		log.Printf("Error Update Bind: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}
	result, err := h.usecase.Update(c.Request().Context(), id, &payload)
	if err != nil {
		log.Printf("Error Update Usecase: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (h *OCRHandler) Detail(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	result, err := h.usecase.Detail(c.Request().Context(), id)
	if err != nil {
		log.Printf("Error Detail Usecase: %v", err)
		return c.JSON(http.StatusNotFound, map[string]any{"message": "not found"})
	}
	return c.JSON(http.StatusOK, result)
}

func (h *OCRHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit, _ = strconv.Atoi(c.QueryParam("offset")) // Fallback jika pakai offset
		if limit == 0 {
			limit = 10
		}
	}

	result, err := h.usecase.List(c.Request().Context(), ocr.ListFilter{
		Vendor: c.QueryParam("vendor"),
		Date:   c.QueryParam("date"),
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		log.Printf("Error List Usecase: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]any{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

func (h *OCRHandler) ExportCSV(c echo.Context) error {
	result, err := h.usecase.List(c.Request().Context(), ocr.ListFilter{Vendor: c.QueryParam("vendor"), Date: c.QueryParam("date"), Page: 1, Limit: 1000})
	if err != nil {
		log.Printf("Error ExportCSV Usecase: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
	}
	c.Response().Header().Set(echo.HeaderContentType, "text/csv")
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=transactions.csv")
	writer := csv.NewWriter(c.Response())
	defer writer.Flush()
	_ = writer.Write([]string{"id", "title", "date_shopping", "amount", "change", "note", "receipt_image", "created_at"})
	for _, row := range result.Data {
		_ = writer.Write([]string{strconv.FormatInt(row.ID, 10), row.Title, row.DateShopping, strconv.FormatFloat(row.Amount, 'f', -1, 64), strconv.FormatInt(row.Change, 10), row.Note, row.ReceiptImage, row.CreatedAt})
	}
	return nil
}
