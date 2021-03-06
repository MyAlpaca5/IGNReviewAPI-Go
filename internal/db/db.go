package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

func NewPool() (*pgxpool.Pool, func(), error) {
	conf, err := pgxpool.ParseConfig(viper.GetString("DATABASE_URL"))
	if err != nil {
		return nil, nil, err
	}

	// TODO: use config file to store those settings
	conf.MaxConns = 25
	conf.MaxConnIdleTime = 15 * time.Minute

	pool, err := pgxpool.ConnectConfig(context.Background(), conf)
	if err != nil {
		return nil, nil, err
	}

	return pool, func() { pool.Close() }, nil
}
