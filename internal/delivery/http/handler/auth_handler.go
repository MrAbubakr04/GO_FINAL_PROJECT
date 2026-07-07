package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	delivery "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/delivery"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/logger"
	authusecase "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/usecase/auth"
)

type AuthHandler struct {
	service *authusecase.Service
}

func NewAuthHandler(service *authusecase.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) CheckPhone(w http.ResponseWriter, r *http.Request) {
	logger.Info("auth check phone route started", logger.Fields{"route": "/auth/check-phone"})
	var req struct {
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("invalid check phone request", logger.Fields{"route": "/auth/check-phone"}, err)
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}

	result, err := h.service.CheckPhone(r.Context(), req.Phone)
	if err != nil {
		logger.Error("check phone failed", logger.Fields{"route": "/auth/check-phone", "phone": req.Phone}, err)
		delivery.WriteError(w, http.StatusInternalServerError, "internal_error", "internal error")
		return
	}

	delivery.WriteSuccess(w, result)
	logger.Info("auth check phone route completed", logger.Fields{"route": "/auth/check-phone"})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	logger.Info("auth register route started", logger.Fields{"route": "/auth/register"})
	var req struct {
		Phone  string `json:"phone"`
		PIN    string `json:"pin"`
		Device string `json:"device"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("invalid register request", logger.Fields{"route": "/auth/register"}, err)
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}

	tokens, err := h.service.Register(r.Context(), authusecase.RegisterInput{Phone: req.Phone, PIN: req.PIN, Device: req.Device})
	if err != nil {
		logger.Error("register failed", logger.Fields{"route": "/auth/register", "phone": req.Phone}, err)
		status := http.StatusBadRequest
		code := "registration_failed"
		if errors.Is(err, domain.ErrInvalidInput) {
			status = http.StatusUnprocessableEntity
			code = "invalid_input"
		} else if errors.Is(err, domain.ErrPhoneAlreadyExists) || errors.Is(err, domain.ErrAccountAlreadyExists) {
			status = http.StatusConflict
			code = "already_exists"
		}
		delivery.WriteError(w, status, code, err.Error())
		return
	}

	delivery.WriteSuccess(w, tokens)
	logger.Info("auth register route completed", logger.Fields{"route": "/auth/register"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger.Info("auth login route started", logger.Fields{"route": "/auth/login"})
	var req struct {
		Phone string `json:"phone"`
		PIN   string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("invalid login request", logger.Fields{"route": "/auth/login"}, err)
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}

	tokens, err := h.service.Login(r.Context(), authusecase.LoginInput{Phone: req.Phone, PIN: req.PIN})
	if err != nil {
		logger.Error("login failed", logger.Fields{"route": "/auth/login", "phone": req.Phone}, err)
		status := http.StatusUnauthorized
		code := "invalid_credentials"
		if errors.Is(err, domain.ErrInvalidInput) {
			status = http.StatusUnprocessableEntity
			code = "invalid_input"
		} else if errors.Is(err, domain.ErrPhoneInactive) {
			status = http.StatusForbidden
			code = "phone_inactive"
		}
		delivery.WriteError(w, status, code, err.Error())
		return
	}

	delivery.WriteSuccess(w, tokens)
	logger.Info("auth login route completed", logger.Fields{"route": "/auth/login"})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	logger.Info("auth refresh route started", logger.Fields{"route": "/auth/refresh"})
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("invalid refresh request", logger.Fields{"route": "/auth/refresh"}, err)
		delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}

	tokens, err := h.service.Refresh(r.Context(), authusecase.RefreshInput{RefreshToken: req.RefreshToken})
	if err != nil {
		logger.Error("refresh failed", logger.Fields{"route": "/auth/refresh"}, err)
		delivery.WriteError(w, http.StatusUnauthorized, "invalid_refresh_token", err.Error())
		return
	}

	delivery.WriteSuccess(w, tokens)
	logger.Info("auth refresh route completed", logger.Fields{"route": "/auth/refresh"})
}
