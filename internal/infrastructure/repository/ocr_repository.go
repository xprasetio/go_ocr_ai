package repository

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"ocr_ai/internal/domain/ocr"
	"ocr_ai/internal/infrastructure/database"
)

type OCRRepository struct{ db *database.D1Client }

func NewOCRRepository(db *database.D1Client) *OCRRepository { return &OCRRepository{db: db} }
func (r *OCRRepository) Migrate(ctx context.Context) error {
	statements := []string{`CREATE TABLE IF NOT EXISTS expenses (id INTEGER PRIMARY KEY AUTOINCREMENT,title TEXT NOT NULL,date_shopping TEXT NULL,"change" INTEGER NULL,amount REAL NULL,note TEXT NULL,receipt_image TEXT NULL,parsed_data TEXT NULL,created_at TEXT DEFAULT CURRENT_TIMESTAMP,updated_at TEXT DEFAULT CURRENT_TIMESTAMP)`, `CREATE TABLE IF NOT EXISTS expense_items (id INTEGER PRIMARY KEY AUTOINCREMENT,expenses_id INTEGER NOT NULL,name TEXT NOT NULL,qty INTEGER NOT NULL,price INTEGER NOT NULL,subtotal INTEGER NOT NULL,created_at TEXT DEFAULT CURRENT_TIMESTAMP,updated_at TEXT DEFAULT CURRENT_TIMESTAMP,FOREIGN KEY (expenses_id) REFERENCES expenses(id) ON DELETE CASCADE)`}
	for _, statement := range statements {
		if _, err := r.db.Exec(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}
func (r *OCRRepository) Create(ctx context.Context, e *ocr.Expense) (*ocr.Expense, error) {
	parsed, _ := json.Marshal(e)
	rows, err := r.db.Exec(ctx, `INSERT INTO expenses (title,date_shopping,"change",amount,note,receipt_image,parsed_data,updated_at) VALUES (?,?,?,?,?,?,?,CURRENT_TIMESTAMP) RETURNING id,created_at,updated_at`, e.Title, e.DateShopping, e.Change, e.Amount, e.Note, e.ReceiptImage, string(parsed))
	if err != nil {
		return nil, err
	}
	e.ID = int64Value(rows[0]["id"])
	e.CreatedAt = stringValue(rows[0]["created_at"])
	e.UpdatedAt = stringValue(rows[0]["updated_at"])
	e.ParsedData = string(parsed)
	for index := range e.Items {
		item := &e.Items[index]
		rows, err := r.db.Exec(ctx, `INSERT INTO expense_items (expenses_id,name,qty,price,subtotal,updated_at) VALUES (?,?,?,?,?,CURRENT_TIMESTAMP) RETURNING id,created_at,updated_at`, e.ID, item.Name, item.Qty, item.Price, item.Subtotal)
		if err != nil {
			return nil, err
		}
		item.ID = int64Value(rows[0]["id"])
		item.ExpensesID = e.ID
		item.CreatedAt = stringValue(rows[0]["created_at"])
		item.UpdatedAt = stringValue(rows[0]["updated_at"])
	}
	return e, nil
}
func (r *OCRRepository) Update(ctx context.Context, id int64, e *ocr.Expense) (*ocr.Expense, error) {
	parsed, _ := json.Marshal(e)
	_, err := r.db.Exec(ctx, `UPDATE expenses SET title=?,date_shopping=?,"change"=?,amount=?,note=?,parsed_data=?,updated_at=CURRENT_TIMESTAMP WHERE id=?`, e.Title, e.DateShopping, e.Change, e.Amount, e.Note, string(parsed), id)
	if err != nil {
		return nil, err
	}
	_, _ = r.db.Exec(ctx, `DELETE FROM expense_items WHERE expenses_id=?`, id)
	for _, item := range e.Items {
		if _, err := r.db.Exec(ctx, `INSERT INTO expense_items (expenses_id,name,qty,price,subtotal,updated_at) VALUES (?,?,?,?,?,CURRENT_TIMESTAMP)`, id, item.Name, item.Qty, item.Price, item.Subtotal); err != nil {
			return nil, err
		}
	}
	return r.FindByID(ctx, id)
}
func (r *OCRRepository) FindByID(ctx context.Context, id int64) (*ocr.Expense, error) {
	rows, err := r.db.Exec(ctx, `SELECT * FROM expenses WHERE id=?`, id)
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	e := expenseFrom(rows[0])
	itemRows, err := r.db.Exec(ctx, `SELECT * FROM expense_items WHERE expenses_id=? ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	for _, row := range itemRows {
		e.Items = append(e.Items, itemFrom(row))
	}
	return &e, nil
}
func (r *OCRRepository) List(ctx context.Context, f ocr.ListFilter) (*ocr.ListResult, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 || f.Limit > 100 {
		f.Limit = 10
	}
	where, params := " WHERE 1=1", []any{}
	if f.Vendor != "" {
		where += " AND title LIKE ?"
		params = append(params, "%"+f.Vendor+"%")
	}
	if f.Date != "" {
		where += " AND date_shopping = ?"
		params = append(params, f.Date)
	}
	countRows, err := r.db.Exec(ctx, `SELECT COUNT(*) total FROM expenses`+where, params...)
	if err != nil {
		return nil, err
	}
	total := int64Value(countRows[0]["total"])
	params = append(params, f.Limit, (f.Page-1)*f.Limit)
	rows, err := r.db.Exec(ctx, `SELECT * FROM expenses`+where+` ORDER BY created_at DESC LIMIT ? OFFSET ?`, params...)
	if err != nil {
		return nil, err
	}
	data := make([]ocr.Expense, 0, len(rows))
	for _, row := range rows {
		data = append(data, expenseFrom(row))
	}
	return &ocr.ListResult{Data: data, Page: f.Page, Limit: f.Limit, Total: total, TotalPages: int(math.Ceil(float64(total) / float64(f.Limit)))}, nil
}
func expenseFrom(row map[string]any) ocr.Expense {
	return ocr.Expense{ID: int64Value(row["id"]), Title: stringValue(row["title"]), DateShopping: stringValue(row["date_shopping"]), Change: int64Value(row["change"]), Amount: floatValue(row["amount"]), Note: stringValue(row["note"]), ReceiptImage: stringValue(row["receipt_image"]), ParsedData: stringValue(row["parsed_data"]), CreatedAt: stringValue(row["created_at"]), UpdatedAt: stringValue(row["updated_at"]), Items: []ocr.ExpenseItem{}}
}
func itemFrom(row map[string]any) ocr.ExpenseItem {
	return ocr.ExpenseItem{ID: int64Value(row["id"]), ExpensesID: int64Value(row["expenses_id"]), Name: stringValue(row["name"]), Qty: int64Value(row["qty"]), Price: int64Value(row["price"]), Subtotal: int64Value(row["subtotal"]), CreatedAt: stringValue(row["created_at"]), UpdatedAt: stringValue(row["updated_at"])}
}
func stringValue(v any) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(v.(string))
}
func int64Value(v any) int64 {
	switch x := v.(type) {
	case float64:
		return int64(x)
	case int64:
		return x
	case int:
		return int64(x)
	case string:
		n, _ := strconv.ParseInt(x, 10, 64)
		return n
	default:
		return 0
	}
}
func floatValue(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int64:
		return float64(x)
	case string:
		n, _ := strconv.ParseFloat(x, 64)
		return n
	default:
		return 0
	}
}
