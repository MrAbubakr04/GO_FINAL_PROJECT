package main

import (
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/logger"
	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/pkg/postgres"
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

	defer db.Close()
}
