package ocr

import (
	"context"
	"fmt"
	"path/filepath"
	"time"
)

type Usecase struct {
	repo     Repository
	analyzer Analyzer
	storage  Storage
}

func NewUsecase(repo Repository, analyzer Analyzer, storage Storage) *Usecase {
	return &Usecase{repo: repo, analyzer: analyzer, storage: storage}
}
func (u *Usecase) Migrate(ctx context.Context) error { return u.repo.Migrate(ctx) }
func (u *Usecase) Analyze(ctx context.Context, fileName, mimeType string, content []byte) (*Expense, error) {
	expense, err := u.analyzer.Analyze(ctx, fileName, mimeType, content)
	if err != nil {
		return nil, err
	}
	storedName := fmt.Sprintf("receipts/%d%s", time.Now().UnixNano(), filepath.Ext(fileName))
	url, err := u.storage.Upload(ctx, storedName, mimeType, content)
	if err != nil {
		return nil, err
	}
	expense.ReceiptImage = url
	return u.repo.Create(ctx, expense)
}
func (u *Usecase) Update(ctx context.Context, id int64, expense *Expense) (*Expense, error) {
	return u.repo.Update(ctx, id, expense)
}
func (u *Usecase) Detail(ctx context.Context, id int64) (*Expense, error) {
	return u.repo.FindByID(ctx, id)
}
func (u *Usecase) List(ctx context.Context, filter ListFilter) (*ListResult, error) {
	return u.repo.List(ctx, filter)
}
