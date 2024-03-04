package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/golang-migrate/migrate/database/clickhouse"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/sixojke/test-service/internal/config"
	delivery "github.com/sixojke/test-service/internal/delivery/http"
	"github.com/sixojke/test-service/internal/repository"
	"github.com/sixojke/test-service/internal/server"
	"github.com/sixojke/test-service/internal/service"
	"github.com/sixojke/test-service/schema"
)

func Run() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("config error: %v", err))
	}

	clickhouse, err := repository.NewClickHouseDB(cfg.ClickHouse)
	if err != nil {
		log.Fatal(fmt.Sprintf("clickhouse connection error: %v", err))
	}
	defer clickhouse.Close()
	log.Info("[CLICKHOUSE] Connection successful")

	if err := schema.MigrateClickHouse(cfg.ClickHouse); err != nil {
		log.Error(fmt.Sprintf("clickhouse migrate error: %v", err))
	}
	log.Info("[CLICKHOUSE] Migrate successful")

	postgres, err := repository.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatal(fmt.Sprintf("postgres connection error: %v", err))
	}
	defer postgres.Close()
	log.Info("[POSTGRES] Connection successful")

	if err := schema.MigratePostgres(cfg.Postgres); err != nil {
		log.Error(fmt.Sprintf("postgres migrate error: %v", err))
	}
	log.Info("[POSTGRES] Migrate successful")

	redis, err := repository.NewRedisDB(cfg.Redis)
	if err != nil {
		log.Fatal(fmt.Sprintf("redis connection error: %v", err))
	}
	defer redis.Close()
	log.Info("[REDIS] Connection successful")

	nats, err := repository.NewNatsClient()
	if err != nil {
		log.Fatal(fmt.Sprintf("nats connection error: %v", err))
	}
	defer nats.Close()
	log.Info("[NATS] Connection successful")
	repo := repository.NewRepository(&repository.Deps{
		ClickHouse: clickhouse,
		Postgres:   postgres,
		Redis:      redis,
		Nats:       nats,
		Config:     cfg,
	})
	services := service.NewService(&service.Deps{
		Config: cfg,
		Repo:   repo,
	})
	services.Goods.GetList(10, 0)
	handler := delivery.NewHandler(services)

	srv := server.NewServer(cfg.HTTPServer, handler.Init())
	go func() {
		if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("error occurred while running http server: %v\n", err)
		}
	}()
	log.Info(fmt.Sprintf("[SERVER] Started :%v", cfg.HTTPServer.Port))

	shutdown(srv, services, clickhouse, postgres, redis, nats)
}

func shutdown(srv *server.Server, services *service.Service, clickhouse, postgres *sqlx.DB, redis *redis.Client, nats *nats.Conn) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 2 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Errorf("failed to stop server: %v", err)
	}

	clickhouse.Close()
	postgres.Close()
	redis.Close()
	nats.Close()

	services.Goods.Shutdown()
}
