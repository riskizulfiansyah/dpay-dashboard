package repository

import (
	"database/sql"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
)

type PaymentRepository interface {
	GetPayments(filter entity.PaymentFilter) ([]entity.Payment, int, error)
	GetStatusCounts(filter entity.PaymentFilter) ([]entity.PaymentStatusCount, error)
}

type Payment struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *Payment {
	return &Payment{db: db}
}

func (r *Payment) buildWhereClause(filter entity.PaymentFilter) (string, []interface{}) {
	where := " WHERE 1=1"
	args := []interface{}{}

	if filter.MerchantID != nil {
		where += " AND merchant_id = ?"
		args = append(args, *filter.MerchantID)
	}

	if filter.Status != nil {
		where += " AND status = ?"
		args = append(args, *filter.Status)
	}

	if filter.FromDate != nil {
		where += " AND created_at >= ?"
		args = append(args, *filter.FromDate)
	}

	if filter.ToDate != nil {
		where += " AND created_at <= ?"
		args = append(args, *filter.ToDate)
	}

	return where, args
}

func (r *Payment) GetPayments(filter entity.PaymentFilter) ([]entity.Payment, int, error) {
	where, args := r.buildWhereClause(filter)

	// Count total matching rows
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM payments" + where
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		return nil, 0, err
	}

	// Fetch paginated rows
	query := `
		SELECT 
			id,
			order_id,
			merchant_name,
			amount,
			status,
			created_at,
			updated_at
		FROM payments` + where

	if filter.SortBy != "" {
		query += " ORDER BY " + filter.SortBy
		if filter.SortDesc {
			query += " DESC"
		} else {
			query += " ASC"
		}
	}

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}
	if filter.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, filter.Offset)
	}

	payments := []entity.Payment{}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment entity.Payment
		if err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.MerchantName,
			&payment.Amount,
			&payment.Status,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		payments = append(payments, payment)
	}

	return payments, totalCount, nil
}

func (r *Payment) GetStatusCounts(filter entity.PaymentFilter) ([]entity.PaymentStatusCount, error) {
	query := `
		SELECT 
			status,
			COUNT(*) as count
		FROM payments
		WHERE 1=1
	`

	args := []interface{}{}

	if filter.FromDate != nil {
		query += " AND created_at >= ?"
		args = append(args, *filter.FromDate)
	}
	if filter.ToDate != nil {
		query += " AND created_at <= ?"
		args = append(args, *filter.ToDate)
	}

	query += " GROUP BY status ORDER BY count DESC"

	var statusCounts []entity.PaymentStatusCount

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sc entity.PaymentStatusCount
		if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
			return nil, err
		}
		statusCounts = append(statusCounts, sc)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return statusCounts, nil
}
