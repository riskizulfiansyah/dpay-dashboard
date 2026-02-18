package handler

import (
	"net/http"

	paymentUsecase "github.com/durianpay/fullstack-boilerplate/internal/module/payment/usecase"
	"github.com/durianpay/fullstack-boilerplate/internal/openapigen"
	"github.com/durianpay/fullstack-boilerplate/internal/transport"
)

type PaymentHandler struct {
	paymentUC paymentUsecase.PaymentUsecase
}

func NewPaymentHandler(paymentUC paymentUsecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUC: paymentUC,
	}
}

func (h *PaymentHandler) GetDashboardV1Payments(w http.ResponseWriter, r *http.Request, params openapigen.GetDashboardV1PaymentsParams) {
	result, err := h.paymentUC.GetPayments(params)
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	response := toPaginatedPaymentResponse(result)

	transport.WriteJSON(w, http.StatusOK, response)
}

func (h *PaymentHandler) GetDashboardV1PaymentsSummary(w http.ResponseWriter, r *http.Request, params openapigen.GetDashboardV1PaymentsSummaryParams) {
	result, err := h.paymentUC.GetPaymentSummary(openapigen.GetDashboardV1PaymentsParams{
		FromDate: params.FromDate,
		ToDate:   params.ToDate,
	})
	if err != nil {
		transport.WriteError(w, err)
		return
	}

	response := toPaymentSummaryResponse(result)

	transport.WriteJSON(w, http.StatusOK, response)
}
