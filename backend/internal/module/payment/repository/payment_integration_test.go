package repository_test

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/payment/repository"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE payments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER NOT NULL,
			merchant_id INTEGER NOT NULL,
			merchant_name TEXT NOT NULL,
			amount INTEGER NOT NULL,
			status TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create payments table: %v", err)
	}

	return db
}

func seedPayments(t *testing.T, db *sql.DB, payments []entity.Payment) {
	t.Helper()
	for _, p := range payments {
		_, err := db.Exec(
			`INSERT INTO payments (order_id, merchant_id, merchant_name, amount, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			p.OrderID, 1, p.MerchantName, p.Amount, p.Status, p.CreatedAt, p.UpdatedAt,
		)
		if err != nil {
			t.Fatalf("failed to seed payment: %v", err)
		}
	}
}

func makeTestPayments(count int, status entity.PaymentStatus, baseTime time.Time) []entity.Payment {
	payments := make([]entity.Payment, count)
	for i := 0; i < count; i++ {
		payments[i] = entity.Payment{
			OrderID:      i + 1,
			MerchantName: "Test Merchant",
			Amount:       (i + 1) * 10000,
			Status:       string(status),
			CreatedAt:    baseTime.Add(time.Duration(i) * time.Hour),
			UpdatedAt:    baseTime.Add(time.Duration(i) * time.Hour),
		}
	}
	return payments
}

func TestGetPayments_Pagination(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	baseTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	seedPayments(t, db, makeTestPayments(10, entity.PaymentStatusCompleted, baseTime))

	repo := repository.NewPaymentRepository(db)

	// Page 1, limit 3
	filter := entity.PaymentFilter{Limit: 3, Offset: 0, SortBy: "id", SortDesc: false}
	results, totalCount, err := repo.GetPayments(filter)
	if err != nil {
		t.Fatalf("GetPayments failed: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 payments, got %d", len(results))
	}
	if totalCount != 10 {
		t.Errorf("expected total count 10, got %d", totalCount)
	}

	// Page 2, limit 3
	filter.Offset = 3
	results, _, err = repo.GetPayments(filter)
	if err != nil {
		t.Fatalf("GetPayments page 2 failed: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 payments on page 2, got %d", len(results))
	}

	// Page 4, limit 3 (only 1 remaining)
	filter.Offset = 9
	results, _, err = repo.GetPayments(filter)
	if err != nil {
		t.Fatalf("GetPayments last page failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 payment on last page, got %d", len(results))
	}
}

func TestGetPayments_DateRange(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	baseTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	seedPayments(t, db, makeTestPayments(10, entity.PaymentStatusCompleted, baseTime))

	repo := repository.NewPaymentRepository(db)

	fromDate := time.Date(2026, 1, 1, 2, 0, 0, 0, time.UTC)
	toDate := time.Date(2026, 1, 1, 5, 0, 0, 0, time.UTC)

	filter := entity.PaymentFilter{
		Limit:    100,
		FromDate: &fromDate,
		ToDate:   &toDate,
		SortBy:   "created_at",
	}

	results, totalCount, err := repo.GetPayments(filter)
	if err != nil {
		t.Fatalf("GetPayments with date range failed: %v", err)
	}

	// Hours 2, 3, 4, 5 → 4 payments
	if len(results) != 4 {
		t.Errorf("expected 4 payments in date range, got %d", len(results))
	}
	if totalCount != 4 {
		t.Errorf("expected total count 4, got %d", totalCount)
	}
}

func TestGetPayments_TotalCount_WithStatusFilter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	baseTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	seedPayments(t, db, makeTestPayments(5, entity.PaymentStatusCompleted, baseTime))
	seedPayments(t, db, makeTestPayments(3, entity.PaymentStatusFailed, baseTime))

	repo := repository.NewPaymentRepository(db)

	// All payments
	_, totalCount, err := repo.GetPayments(entity.PaymentFilter{Limit: 100})
	if err != nil {
		t.Fatalf("GetPayments failed: %v", err)
	}
	if totalCount != 8 {
		t.Errorf("expected total count 8, got %d", totalCount)
	}

	// Only completed
	status := string(entity.PaymentStatusCompleted)
	_, totalCount, err = repo.GetPayments(entity.PaymentFilter{Limit: 100, Status: &status})
	if err != nil {
		t.Fatalf("GetPayments with status failed: %v", err)
	}
	if totalCount != 5 {
		t.Errorf("expected 5 completed payments, got %d", totalCount)
	}
}

func TestGetStatusCounts(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	baseTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	seedPayments(t, db, makeTestPayments(5, entity.PaymentStatusCompleted, baseTime))
	seedPayments(t, db, makeTestPayments(3, entity.PaymentStatusFailed, baseTime))
	seedPayments(t, db, makeTestPayments(2, entity.PaymentStatusProcessing, baseTime))

	repo := repository.NewPaymentRepository(db)

	counts, err := repo.GetStatusCounts(entity.PaymentFilter{})
	if err != nil {
		t.Fatalf("GetStatusCounts failed: %v", err)
	}

	statusMap := make(map[string]int)
	for _, c := range counts {
		statusMap[c.Status] = c.Count
	}

	if statusMap[string(entity.PaymentStatusCompleted)] != 5 {
		t.Errorf("expected 5 completed, got %d", statusMap[string(entity.PaymentStatusCompleted)])
	}
	if statusMap[string(entity.PaymentStatusFailed)] != 3 {
		t.Errorf("expected 3 failed, got %d", statusMap[string(entity.PaymentStatusFailed)])
	}
	if statusMap[string(entity.PaymentStatusProcessing)] != 2 {
		t.Errorf("expected 2 processing, got %d", statusMap[string(entity.PaymentStatusProcessing)])
	}
}

func TestGetStatusCounts_WithDateFilter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	baseTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	seedPayments(t, db, makeTestPayments(5, entity.PaymentStatusCompleted, baseTime))
	seedPayments(t, db, makeTestPayments(3, entity.PaymentStatusFailed, baseTime))

	repo := repository.NewPaymentRepository(db)

	fromDate := time.Date(2026, 1, 1, 3, 0, 0, 0, time.UTC)
	counts, err := repo.GetStatusCounts(entity.PaymentFilter{FromDate: &fromDate})
	if err != nil {
		t.Fatalf("GetStatusCounts with date filter failed: %v", err)
	}

	statusMap := make(map[string]int)
	for _, c := range counts {
		statusMap[c.Status] = c.Count
	}

	// completed: hours 0,1,2,3,4 → from hour 3: 2 payments
	if statusMap[string(entity.PaymentStatusCompleted)] != 2 {
		t.Errorf("expected 2 completed from date, got %d", statusMap[string(entity.PaymentStatusCompleted)])
	}
}
