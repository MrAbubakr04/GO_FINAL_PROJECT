package http

import (
	"net/http"

	authhandler "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/delivery/http/handler"
	authusecase "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/usecase/auth"
)

func NewRouter(authService *authusecase.Service) http.Handler {
	mux := http.NewServeMux()
	authHandler := authhandler.NewAuthHandler(authService)

	mux.HandleFunc("/auth/check-phone", authHandler.CheckPhone)
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/auth/refresh", authHandler.Refresh)

	return mux
}
