package service

import (
	"github.com/sixojke/test-service/internal/config"
	"github.com/sixojke/test-service/internal/domain"
	"github.com/sixojke/test-service/internal/repository"
)

type Goods interface {
	Create(inp domain.ItemCreateInp) (*domain.Item, error)
	Update(inp domain.ItemUpdateInp) (*domain.Item, error)
	Delete(inp domain.ItemDeleteInp) (*domain.ItemDeleteOut, error)
	GetList(limit, offset int) (*domain.List, error)
	Reprioritiize(inp domain.ItemReprioritiizeInp) ([]domain.ItemReprioritiizeOut, error)
	Shutdown()
}

type Deps struct {
	Config *config.Config
	Repo   *repository.Repository
}

type Service struct {
	Goods Goods
}

func NewService(deps *Deps) *Service {
	return &Service{
		Goods: NewGoodsService(deps.Repo.Goods, deps.Repo.Log, deps.Repo.Cache, deps.Config.Service),
	}
}
