package handler

import (
	"fmt"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	paymentUsecase "github.com/durianpay/fullstack-boilerplate/internal/module/payment/usecase"
	"github.com/durianpay/fullstack-boilerplate/internal/openapigen"
)

func toPaginatedPaymentResponse(result *paymentUsecase.PaginatedPayments) openapigen.PaymentListResponse {
	payments := toPaymentResponse(result.Payments)
	page := result.Page
	limit := result.Limit
	totalCount := result.TotalCount
	totalPages := result.TotalPages

	return openapigen.PaymentListResponse{
		Payments: &payments,
		Pagination: &openapigen.Pagination{
			Page:       &page,
			Limit:      &limit,
			TotalCount: &totalCount,
			TotalPages: &totalPages,
		},
	}
}

func toPaymentResponse(payments []entity.Payment) []openapigen.Payment {
	responsePayments := make([]openapigen.Payment, 0, len(payments))
	for _, p := range payments {
		amountStr := fmt.Sprintf("%d", p.Amount)
		idStr := fmt.Sprintf("%d", p.ID)
		merchantName := p.MerchantName
		status := p.Status

		responsePayments = append(responsePayments, openapigen.Payment{
			Id:        &idStr,
			Amount:    &amountStr,
			Merchant:  &merchantName,
			Status:    &status,
			CreatedAt: &p.CreatedAt,
		})
	}
	return responsePayments
}

func toPaymentSummaryResponse(summary *entity.PaymentSummary) openapigen.PaymentSummaryResponse {
	statusCounts := make([]openapigen.PaymentStatusCount, 0, len(summary.StatusCounts))
	for _, sc := range summary.StatusCounts {
		status := sc.Status
		count := sc.Count
		statusCounts = append(statusCounts, openapigen.PaymentStatusCount{
			Status: &status,
			Count:  &count,
		})
	}

	total := summary.Total
	return openapigen.PaymentSummaryResponse{
		StatusCounts: &statusCounts,
		Total:        &total,
	}
}
