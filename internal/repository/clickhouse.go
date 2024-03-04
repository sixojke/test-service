package repository

import (
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/sixojke/test-service/internal/config"
)

func NewClickHouseDB(cfg config.ClickHouseConfig) (*sqlx.DB, error) {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{cfg.Host + ":" + cfg.Port},
		Auth: clickhouse.Auth{
			Database: cfg.DBName,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		Debug: true,
	})

	db := sqlx.NewDb(conn, "clickhouse")

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %v", err)
	}

	return db, nil
}
