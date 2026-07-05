package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConfig struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DbName string

	MaxConns int32
	MinConns int32
}

func New(db DbConfig) (*pgxpool.Pool, error) {
	db_url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		db.User,
		db.Pass,
		db.Host,
		db.Port,
		db.DbName,
	)

	poolConf, err := pgxpool.ParseConfig(db_url)
	if err != nil {
		return nil, err
	}

	poolConf.MaxConns = db.MaxConns
	poolConf.MinConns = db.MinConns
	poolConf.MaxConnLifetime = time.Hour
	poolConf.MaxConnIdleTime = 15 * time.Minute
	poolConf.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConf)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, err
}
