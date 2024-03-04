package repository

import (
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/sixojke/test-service/internal/config"
	"github.com/sixojke/test-service/internal/domain"
)

type Goods interface {
	Create(inp domain.ItemCreateInp) (*domain.Item, error)
	Update(inp domain.ItemUpdateInp) (*domain.Item, error)
	Delete(inp domain.ItemDeleteInp) (*domain.ItemDeleteOut, error)
	GetList(limit, offset int) (*domain.List, error)
	GetById(id int) (*domain.Item, error)
	Reprioritiize(domain.ItemReprioritiizeInp) ([]domain.ItemReprioritiizeOut, error)
}

type Log interface {
	Send(msg []*domain.ChangesHistory) error
}

type Cache interface {
	SetList(limit, offset int, list *domain.List) error
	GetList(limit, offset int) (*domain.List, error)
	Delete(itemId int) error
}

type Deps struct {
	ClickHouse *sqlx.DB
	Postgres   *sqlx.DB
	Redis      *redis.Client
	Nats       *nats.Conn
	Config     *config.Config
}

type Repository struct {
	Goods Goods
	Log   Log
	Cache Cache
}

func NewRepository(deps *Deps) *Repository {
	return &Repository{
		Goods: NewGoodsPostgres(deps.Postgres),
		Log:   NewLogNats(deps.Nats),
		Cache: NewCacheRedis(deps.Redis, deps.Config.Cache),
	}
}
