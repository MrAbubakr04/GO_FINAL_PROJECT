package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	delivery "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/delivery"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/logger"
	clientusecase "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/usecase/client"
)

type ClientHandler struct {
	service *clientusecase.Service
}

func NewClientHandler(service *clientusecase.Service) *ClientHandler {
	return &ClientHandler{service: service}
}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	logger.Info("client create route started", logger.Fields{"route": "/clients"})
	var req clientusecase.CreateClientInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}
	client, err := h.service.CreateClient(r.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		code := "create_failed"
		if errors.Is(err, domain.ErrInvalidInput) {
			status = http.StatusUnprocessableEntity
			code = "invalid_input"
		} else if errors.Is(err, domain.ErrClientAlreadyExists) {
			status = http.StatusConflict
			code = "client_already_exists"
		}
		delivery.WriteError(w, status, code, err.Error())
		return
	}
	delivery.WriteSuccess(w, client)
}

func (h *ClientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}
	client, err := h.service.GetClientByID(r.Context(), id)
	if err != nil {
		delivery.WriteError(w, http.StatusNotFound, "client_not_found", err.Error())
		return
	}
	delivery.WriteSuccess(w, client)
}

func (h *ClientHandler) GetByTIN(w http.ResponseWriter, r *http.Request) {
	tin := r.URL.Query().Get("tin")
	client, err := h.service.GetClientByTIN(r.Context(), tin)
	if err != nil {
		delivery.WriteError(w, http.StatusNotFound, "client_not_found", err.Error())
		return
	}
	delivery.WriteSuccess(w, client)
}

func (h *ClientHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}
	var req clientusecase.UpdateClientInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}
	client, err := h.service.UpdateClient(r.Context(), id, req)
	if err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "update_failed", err.Error())
		return
	}
	delivery.WriteSuccess(w, client)
}

func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}
	if err := h.service.DeactivateClient(r.Context(), id); err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "deactivate_failed", err.Error())
		return
	}
	delivery.WriteSuccess(w, map[string]bool{"deleted": true})
}

func (h *ClientHandler) Identify(w http.ResponseWriter, r *http.Request) {
	var req clientusecase.IdentifyInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}
	ok, err := h.service.IdentifyClient(r.Context(), req)
	if err != nil {
		delivery.WriteError(w, http.StatusBadRequest, "identification_failed", err.Error())
		return
	}
	delivery.WriteSuccess(w, map[string]bool{"identified": ok})
}
