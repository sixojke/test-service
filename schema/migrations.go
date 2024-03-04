package schema

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	"github.com/sixojke/test-service/internal/config"
	"github.com/sixojke/test-service/pkg/utils"
)

func MigrateClickHouse(cfg config.ClickHouseConfig) error {
	err := migrateUp(fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName), "/clickhouse")
	if err != nil {
		return err
	}

	return nil
}

func MigratePostgres(cfg config.PostgresConfig) error {
	err := migrateUp(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode), "/postgres")
	if err != nil {
		return err
	}

	return nil
}

func migrateUp(conn, dir string) error {
	path, err := utils.CustomPath("/schema")
	if err != nil {
		return err
	}
	m, err := migrate.New(
		"file:"+path+dir,
		conn,
	)
	if err != nil {
		return fmt.Errorf("create migration: %s", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("load migration: %s", err)
		}
	}

	return nil
}
