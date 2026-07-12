package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	delivery "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/delivery"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/logger"
	transactionusecase "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/usecase/transaction"
)

type TransactionHandler struct {
	service *transactionusecase.Service
}

func NewTransactionHandler(service *transactionusecase.Service) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	logger.Info("transaction deposit route started", logger.Fields{"route": "/transactions/deposit"})
	var req transactionusecase.DepositInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}
	result, err := h.service.Deposit(r.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		code := "deposit_failed"
		switch {
		case errors.Is(err, domain.ErrPhoneNotFound):
			status = http.StatusNotFound
			code = "phone_not_found"
		case errors.Is(err, domain.ErrPhoneInactive):
			status = http.StatusForbidden
			code = "phone_inactive"
		case errors.Is(err, domain.ErrAccountNotFound):
			status = http.StatusNotFound
			code = "account_not_found"
		case errors.Is(err, domain.ErrAccountInactive):
			status = http.StatusForbidden
			code = "account_inactive"
		case errors.Is(err, domain.ErrInvalidAmount):
			status = http.StatusUnprocessableEntity
			code = "invalid_amount"
		}
		delivery.WriteError(w, status, code, err.Error())
		return
	}
	delivery.WriteSuccess(w, result)
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	logger.Info("transaction transfer route started", logger.Fields{"route": "/transactions/transfer"})
	var req transactionusecase.TransferInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}
	result, err := h.service.Transfer(r.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		code := "transfer_failed"
		switch {
		case errors.Is(err, domain.ErrSameAccountTransfer):
			status = http.StatusConflict
			code = "same_account_transfer"
		case errors.Is(err, domain.ErrInvalidAmount):
			status = http.StatusUnprocessableEntity
			code = "invalid_amount"
		case errors.Is(err, domain.ErrInvalidInput):
			status = http.StatusUnprocessableEntity
			code = "invalid_input"
		case errors.Is(err, domain.ErrPhoneNotFound):
			status = http.StatusNotFound
			code = "phone_not_found"
		case errors.Is(err, domain.ErrPhoneInactive):
			status = http.StatusForbidden
			code = "phone_inactive"
		case errors.Is(err, domain.ErrAccountNotFound):
			status = http.StatusNotFound
			code = "account_not_found"
		case errors.Is(err, domain.ErrAccountInactive):
			status = http.StatusForbidden
			code = "account_inactive"
		case errors.Is(err, domain.ErrAccountNotIdentified):
			status = http.StatusForbidden
			code = "account_not_identified"
		case errors.Is(err, domain.ErrInsufficientFunds):
			status = http.StatusConflict
			code = "insufficient_funds"
		}
		delivery.WriteError(w, status, code, err.Error())
		return
	}
	delivery.WriteSuccess(w, result)
}

func (h *TransactionHandler) AccountHistory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("account_id"), 10, 64)
	if err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}
	history, err := h.service.AccountHistory(r.Context(), id)
	if err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "history_failed", err.Error())
		return
	}
	delivery.WriteSuccess(w, history)
}

func (h *TransactionHandler) ClientHistory(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("client_id"), 10, 64)
	if err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}
	history, err := h.service.ClientHistory(r.Context(), id)
	if err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "history_failed", err.Error())
		return
	}
	delivery.WriteSuccess(w, history)
}
