package ocr

import "context"

type Analyzer interface {
	Analyze(ctx context.Context, fileName string, mimeType string, content []byte) (*Expense, error)
}

type Storage interface {
	Upload(ctx context.Context, fileName string, mimeType string, content []byte) (string, error)
}
