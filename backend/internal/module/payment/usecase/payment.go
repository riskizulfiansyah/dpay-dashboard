package usecase

import (
	"math"
	"strings"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/payment/repository"
	"github.com/durianpay/fullstack-boilerplate/internal/openapigen"
)

type PaginatedPayments struct {
	Payments   []entity.Payment
	Page       int
	Limit      int
	TotalCount int
	TotalPages int
}

type PaymentUsecase interface {
	GetPayments(params openapigen.GetDashboardV1PaymentsParams) (*PaginatedPayments, error)
	GetPaymentSummary(params openapigen.GetDashboardV1PaymentsParams) (*entity.PaymentSummary, error)
}

type Payment struct {
	repo repository.PaymentRepository
}

func NewPaymentUsecase(repo repository.PaymentRepository) *Payment {
	return &Payment{repo: repo}
}

var (
	validStatuses = entity.ValidPaymentStatuses
	allowedSorts  = map[string]bool{
		"id":         true,
		"amount":     true,
		"status":     true,
		"created_at": true,
		"updated_at": true,
	}
)

const (
	defaultPage  = 1
	defaultLimit = 20
	maxLimit     = 100
)

func (p *Payment) GetPayments(params openapigen.GetDashboardV1PaymentsParams) (*PaginatedPayments, error) {
	page := defaultPage
	limit := defaultLimit

	if params.Page != nil && *params.Page > 0 {
		page = *params.Page
	}
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
		if limit > maxLimit {
			limit = maxLimit
		}
	}

	filter := entity.PaymentFilter{
		Limit:    limit,
		Offset:   (page - 1) * limit,
		FromDate: params.FromDate,
		ToDate:   params.ToDate,
	}

	if params.Status != nil {
		if !validStatuses[entity.PaymentStatus(*params.Status)] {
			return nil, entity.ErrorBadRequest("invalid status")
		}
		filter.Status = params.Status
	}

	if params.Sort != nil {
		sortField := *params.Sort
		if strings.HasPrefix(sortField, "-") {
			filter.SortBy = strings.TrimPrefix(sortField, "-")
			filter.SortDesc = true
		} else {
			filter.SortBy = sortField
		}
		if !allowedSorts[filter.SortBy] {
			return nil, entity.ErrorBadRequest("invalid sort field")
		}
	} else {
		filter.SortBy = "created_at"
		filter.SortDesc = true
	}

	payments, totalCount, err := p.repo.GetPayments(filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	return &PaginatedPayments{
		Payments:   payments,
		Page:       page,
		Limit:      limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (p *Payment) GetPaymentSummary(params openapigen.GetDashboardV1PaymentsParams) (*entity.PaymentSummary, error) {
	filter := entity.PaymentFilter{
		FromDate: params.FromDate,
		ToDate:   params.ToDate,
	}

	counts, err := p.repo.GetStatusCounts(filter)
	if err != nil {
		return nil, err
	}

	statusMap := make(map[string]int)
	for _, c := range counts {
		statusMap[c.Status] = c.Count
	}

	allStatuses := entity.AllPaymentStatuses
	finalCounts := []entity.PaymentStatusCount{}
	total := 0

	for _, s := range allStatuses {
		count := statusMap[string(s)]
		finalCounts = append(finalCounts, entity.PaymentStatusCount{
			Status: string(s),
			Count:  count,
		})
		total += count
	}

	return &entity.PaymentSummary{
		StatusCounts: finalCounts,
		Total:        total,
	}, nil
}
