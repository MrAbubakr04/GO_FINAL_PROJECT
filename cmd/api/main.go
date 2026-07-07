package main

import (
	"log"
	"net/http"

	internalhttp "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/delivery/http"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/logger"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/pkg/postgres"
	repositorypostgres "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/repository/postgres"
	authusecase "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/usecase/auth"
	clientusecase "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/usecase/client"
	transactionusecase "github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/usecase/transaction"
)

func main() {
	if err := logger.Init(&logger.Options{Directory: "logs", FileName: "app.jsonl"}); err != nil {
		panic(err)
	}
	logger.Debug("Application started compleat", nil)

	db, err := postgres.New(postgres.DbConfig{Host: "localhost", Port: 5432, User: "postgres", Pass: "postgres", DbName: "humo_online", MaxConns: 10, MinConns: 2})
	if err != nil {
		logger.Error("Соединение с БД: Postgres не установлено.", nil, err)
		panic(err)
	}
	logger.Debug("Соединение с БД: Postgres успешно установлено.", nil)

	authRepo := repositorypostgres.NewAuthRepository(db)
	clientRepo := repositorypostgres.NewClientRepository(db)
	transactionRepo := repositorypostgres.NewTransactionRepository(db)
	authService := authusecase.NewService(authRepo, nil)
	clientService := clientusecase.NewService(clientRepo)
	transactionService := transactionusecase.NewService(transactionRepo)
	router := internalhttp.NewRouter(authService, clientService, transactionService)

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("server stopped", nil, err)
		panic(err)
	}

	defer db.Close()
}
