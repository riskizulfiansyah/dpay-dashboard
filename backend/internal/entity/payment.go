package entity

import "time"

type PaymentStatus string

const (
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusFailed     PaymentStatus = "failed"
)

var (
	AllPaymentStatuses = []PaymentStatus{
		PaymentStatusCompleted,
		PaymentStatusProcessing,
		PaymentStatusFailed,
	}

	ValidPaymentStatuses = map[PaymentStatus]bool{
		PaymentStatusCompleted:  true,
		PaymentStatusProcessing: true,
		PaymentStatusFailed:     true,
	}
)

type Payment struct {
	ID           int
	OrderID      int
	MerchantName string
	Amount       int
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PaymentFilter struct {
	MerchantID *int
	Status     *string
	Limit      int
	Offset     int
	FromDate   *time.Time
	ToDate     *time.Time
	SortBy     string
	SortDesc   bool
}

type PaymentStatusCount struct {
	Status string
	Count  int
}

type PaymentSummary struct {
	StatusCounts []PaymentStatusCount
	Total        int
}
