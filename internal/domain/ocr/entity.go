package ocr

type Expense struct {
	ID           int64         `json:"id"`
	Title        string        `json:"title"`
	DateShopping string        `json:"date_shopping"`
	Change       int64         `json:"change"`
	Amount       float64       `json:"amount"`
	Note         string        `json:"note"`
	ReceiptImage string        `json:"receipt_image"`
	ParsedData   string        `json:"parsed_data"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
	Items        []ExpenseItem `json:"items"`
}

type ExpenseItem struct {
	ID         int64  `json:"id"`
	ExpensesID int64  `json:"expenses_id"`
	Name       string `json:"name"`
	Qty        int64  `json:"qty"`
	Price      int64  `json:"price"`
	Subtotal   int64  `json:"subtotal"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type ListFilter struct {
	Vendor string
	Date   string
	Page   int
	Limit  int
}

type ListResult struct {
	Data       []Expense `json:"data"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	Total      int64     `json:"total"`
	TotalPages int       `json:"total_pages"`
}
