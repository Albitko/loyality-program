package entities

type Order struct {
	OrderID string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type OrderWithTime struct {
	OrderID   string  `json:"number"`
	Status    string  `json:"status"`
	Accrual   float64 `json:"accrual,omitempty"`
	UpdatedAt string  `json:"uploaded_at"`
}
