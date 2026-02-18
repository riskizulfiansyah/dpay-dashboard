package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/payment/usecase"
	"github.com/durianpay/fullstack-boilerplate/internal/openapigen"
)

type mockPaymentRepo struct {
	payments     []entity.Payment
	totalCount   int
	statusCounts []entity.PaymentStatusCount
	lastFilter   entity.PaymentFilter
	err          error
}

func (m *mockPaymentRepo) GetPayments(filter entity.PaymentFilter) ([]entity.Payment, int, error) {
	m.lastFilter = filter
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.payments, m.totalCount, nil
}

func (m *mockPaymentRepo) GetStatusCounts(filter entity.PaymentFilter) ([]entity.PaymentStatusCount, error) {
	m.lastFilter = filter
	if m.err != nil {
		return nil, m.err
	}
	return m.statusCounts, nil
}

func TestGetPayments_DefaultPagination(t *testing.T) {
	mock := &mockPaymentRepo{
		payments:   []entity.Payment{{ID: 1}},
		totalCount: 1,
	}
	uc := usecase.NewPaymentUsecase(mock)

	result, err := uc.GetPayments(openapigen.GetDashboardV1PaymentsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Page != 1 {
		t.Errorf("expected page 1, got %d", result.Page)
	}
	if result.Limit != 20 {
		t.Errorf("expected limit 20, got %d", result.Limit)
	}
	if mock.lastFilter.Offset != 0 {
		t.Errorf("expected offset 0, got %d", mock.lastFilter.Offset)
	}
}

func TestGetPayments_CustomPagination(t *testing.T) {
	mock := &mockPaymentRepo{
		payments:   []entity.Payment{},
		totalCount: 50,
	}
	uc := usecase.NewPaymentUsecase(mock)

	page := 3
	limit := 10
	result, err := uc.GetPayments(openapigen.GetDashboardV1PaymentsParams{
		Page:  &page,
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Page != 3 {
		t.Errorf("expected page 3, got %d", result.Page)
	}
	if result.Limit != 10 {
		t.Errorf("expected limit 10, got %d", result.Limit)
	}
	if result.TotalCount != 50 {
		t.Errorf("expected total count 50, got %d", result.TotalCount)
	}
	if result.TotalPages != 5 {
		t.Errorf("expected 5 total pages, got %d", result.TotalPages)
	}
	if mock.lastFilter.Offset != 20 {
		t.Errorf("expected offset 20, got %d", mock.lastFilter.Offset)
	}
}

func TestGetPayments_LimitCapped(t *testing.T) {
	mock := &mockPaymentRepo{
		payments:   []entity.Payment{},
		totalCount: 0,
	}
	uc := usecase.NewPaymentUsecase(mock)

	limit := 500
	result, err := uc.GetPayments(openapigen.GetDashboardV1PaymentsParams{
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Limit != 100 {
		t.Errorf("expected limit capped at 100, got %d", result.Limit)
	}
}

func TestGetPayments_DateRange(t *testing.T) {
	mock := &mockPaymentRepo{
		payments:   []entity.Payment{},
		totalCount: 0,
	}
	uc := usecase.NewPaymentUsecase(mock)

	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)

	_, err := uc.GetPayments(openapigen.GetDashboardV1PaymentsParams{
		FromDate: &from,
		ToDate:   &to,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mock.lastFilter.FromDate == nil || !mock.lastFilter.FromDate.Equal(from) {
		t.Errorf("expected FromDate %v, got %v", from, mock.lastFilter.FromDate)
	}
	if mock.lastFilter.ToDate == nil || !mock.lastFilter.ToDate.Equal(to) {
		t.Errorf("expected ToDate %v, got %v", to, mock.lastFilter.ToDate)
	}
}

func TestGetPayments_InvalidStatus(t *testing.T) {
	mock := &mockPaymentRepo{}
	uc := usecase.NewPaymentUsecase(mock)

	badStatus := "invalid_status"
	_, err := uc.GetPayments(openapigen.GetDashboardV1PaymentsParams{
		Status: &badStatus,
	})
	if err == nil {
		t.Fatal("expected error for invalid status, got nil")
	}
}

func TestGetPayments_TotalPages_Rounding(t *testing.T) {
	mock := &mockPaymentRepo{
		payments:   []entity.Payment{},
		totalCount: 7,
	}
	uc := usecase.NewPaymentUsecase(mock)

	limit := 3
	result, err := uc.GetPayments(openapigen.GetDashboardV1PaymentsParams{
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 7 / 3 = 2.33 → 3 pages
	if result.TotalPages != 3 {
		t.Errorf("expected 3 total pages (ceil), got %d", result.TotalPages)
	}
}

func TestGetPayments_SortDescending(t *testing.T) {
	mock := &mockPaymentRepo{
		payments:   []entity.Payment{},
		totalCount: 0,
	}
	uc := usecase.NewPaymentUsecase(mock)

	sort := "-amount"
	_, err := uc.GetPayments(openapigen.GetDashboardV1PaymentsParams{
		Sort: &sort,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mock.lastFilter.SortDesc {
		t.Error("expected SortDesc to be true")
	}
	if mock.lastFilter.SortBy != "amount" {
		t.Errorf("expected SortBy 'amount', got '%s'", mock.lastFilter.SortBy)
	}
}

func TestGetPayments_InvalidSortField(t *testing.T) {
	mock := &mockPaymentRepo{}
	uc := usecase.NewPaymentUsecase(mock)

	sort := "invalid_field"
	_, err := uc.GetPayments(openapigen.GetDashboardV1PaymentsParams{
		Sort: &sort,
	})
	if err == nil {
		t.Fatal("expected error for invalid sort field, got nil")
	}
}

func TestGetPaymentSummary_AllStatuses(t *testing.T) {
	mock := &mockPaymentRepo{
		statusCounts: []entity.PaymentStatusCount{
			{Status: string(entity.PaymentStatusCompleted), Count: 10},
			{Status: string(entity.PaymentStatusProcessing), Count: 5},
			{Status: string(entity.PaymentStatusFailed), Count: 2},
		},
	}
	uc := usecase.NewPaymentUsecase(mock)

	summary, err := uc.GetPaymentSummary(openapigen.GetDashboardV1PaymentsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.Total != 17 {
		t.Errorf("expected total 17, got %d", summary.Total)
	}
	if len(summary.StatusCounts) != 3 {
		t.Errorf("expected 3 status counts, got %d", len(summary.StatusCounts))
	}
}

func TestGetPaymentSummary_MissingStatus(t *testing.T) {
	mock := &mockPaymentRepo{
		statusCounts: []entity.PaymentStatusCount{
			{Status: string(entity.PaymentStatusCompleted), Count: 10},
		},
	}
	uc := usecase.NewPaymentUsecase(mock)

	summary, err := uc.GetPaymentSummary(openapigen.GetDashboardV1PaymentsParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 10 + 0 + 0 = 10
	if summary.Total != 10 {
		t.Errorf("expected total 10, got %d", summary.Total)
	}

	// Should still have entries for all statuses, even if count is 0
	if len(summary.StatusCounts) != 3 {
		t.Errorf("expected 3 status counts, got %d", len(summary.StatusCounts))
	}
}

func TestGetPaymentSummary_RepoError(t *testing.T) {
	mock := &mockPaymentRepo{
		err: errors.New("db error"),
	}
	uc := usecase.NewPaymentUsecase(mock)

	_, err := uc.GetPaymentSummary(openapigen.GetDashboardV1PaymentsParams{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
