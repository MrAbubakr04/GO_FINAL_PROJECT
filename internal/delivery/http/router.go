package http

import (
	"encoding/json"
	"net/http"

	delivery "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/delivery"
	authhandler "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/delivery/http/handler"
	authusecase "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/usecase/auth"
	clientusecase "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/usecase/client"
)

func NewRouter(authService *authusecase.Service, clientService *clientusecase.Service) http.Handler {
	mux := http.NewServeMux()
	authHandler := authhandler.NewAuthHandler(authService)
	clientHandler := authhandler.NewClientHandler(clientService)

	mux.HandleFunc("/auth/check-phone", authHandler.CheckPhone)
	mux.HandleFunc("/auth/register", authHandler.Register)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/auth/refresh", authHandler.Refresh)
	mux.HandleFunc("/employees/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			delivery.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
			return
		}
		user, err := clientService.Login(r.Context(), clientusecase.LoginInput{Login: req.Login, Password: req.Password})
		if err != nil {
			delivery.WriteError(w, http.StatusUnauthorized, "invalid_credentials", err.Error())
			return
		}
		delivery.WriteSuccess(w, user)
	})
	mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			clientHandler.Create(w, r)
		case http.MethodGet:
			clientHandler.GetByTIN(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/clients/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet:
			clientHandler.GetByID(w, r)
		case r.Method == http.MethodPut:
			clientHandler.Update(w, r)
		case r.Method == http.MethodDelete:
			clientHandler.Delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/clients/identify", clientHandler.Identify)

	return mux
}
