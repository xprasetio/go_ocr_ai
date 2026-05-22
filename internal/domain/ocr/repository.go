package ocr

import "context"

type Repository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, expense *Expense) (*Expense, error)
	Update(ctx context.Context, id int64, expense *Expense) (*Expense, error)
	FindByID(ctx context.Context, id int64) (*Expense, error)
	List(ctx context.Context, filter ListFilter) (*ListResult, error)
}
