package service

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/sixojke/test-service/internal/config"
	"github.com/sixojke/test-service/internal/domain"
	"github.com/sixojke/test-service/internal/repository"
)

type GoodsService struct {
	repo   repository.Goods
	log    repository.Log
	cache  repository.Cache
	config config.ServiceConfig
}

func NewGoodsService(repo repository.Goods, log repository.Log, cache repository.Cache, cfg config.ServiceConfig) *GoodsService {
	return &GoodsService{
		repo:   repo,
		log:    log,
		cache:  cache,
		config: cfg,
	}
}

func (s *GoodsService) Create(inp domain.ItemCreateInp) (*domain.Item, error) {
	if err := inp.Validate(); err != nil {
		return nil, fmt.Errorf("create: %v", err)
	}

	item, err := s.repo.Create(inp)
	if err != nil {
		return nil, fmt.Errorf("create: %v", err)
	}

	go s.writeLog(item.Id)

	return item, nil
}

func (s *GoodsService) Update(inp domain.ItemUpdateInp) (*domain.Item, error) {
	if err := s.cache.Delete(inp.Id); err != nil {
		log.Error(fmt.Errorf("[CACHE] ERROR!!! | %v", err))
		return nil, fmt.Errorf("update: %v", err)
	}

	item, err := s.repo.Update(inp)
	if err != nil {
		return nil, fmt.Errorf("update: %v", err)
	}

	go s.writeLog(item.Id)

	return item, nil
}

func (s *GoodsService) Delete(inp domain.ItemDeleteInp) (*domain.ItemDeleteOut, error) {
	if err := s.cache.Delete(inp.Id); err != nil {
		log.Error(fmt.Errorf("[CACHE] ERROR!!! | %v", err))
		return nil, fmt.Errorf("delete: %v", err)
	}

	item, err := s.repo.Delete(inp)
	if err != nil {
		return nil, fmt.Errorf("delete: %v", err)
	}

	go s.writeLog(item.Id)

	return item, err
}

func (s *GoodsService) GetList(limit, offset int) (*domain.List, error) {
	if limit > s.config.MaxLimit {
		limit = s.config.MaxLimit
	}

	return s.getList(limit, offset)
}

func (s *GoodsService) getList(limit, offset int) (*domain.List, error) {
	listRedis, err := s.cache.GetList(limit, offset)
	if err != nil {
		log.Errorf("[CACHE] ERROR!!! | %v", err)
	}
	log.Info(listRedis)

	if listRedis == nil {
		return s.getListPostgres(limit, offset)
	}

	return listRedis, err
}

func (s *GoodsService) getListPostgres(limit, offset int) (*domain.List, error) {
	listPostgres, err := s.repo.GetList(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get list: %v", err)
	}

	if err := s.cache.SetList(limit, offset, listPostgres); err != nil {
		return nil, fmt.Errorf("get list: %v", err)
	}

	return listPostgres, nil
}

func (s *GoodsService) Reprioritiize(inp domain.ItemReprioritiizeInp) ([]domain.ItemReprioritiizeOut, error) {
	goods, err := s.repo.Reprioritiize(inp)
	if err != nil {
		return nil, fmt.Errorf("reprioritiize: %v", err)
	}

	for _, item := range goods {
		if err := s.cache.Delete(item.Id); err != nil {
			return nil, fmt.Errorf("reprioritiize: %v", err)
		}
		s.writeLog(item.Id)
	}

	return goods, nil
}

var logs = make([]*domain.ChangesHistory, 0, 50)

func (s *GoodsService) writeLog(itemId int) {
	logItem, err := s.repo.GetById(itemId)
	if err != nil {
		log.Warn(fmt.Errorf("error: %v", err))
	}
	if logItem != nil {
		logs = append(logs, &domain.ChangesHistory{
			Id:          logItem.Id,
			ProjectId:   logItem.ProjectId,
			Name:        logItem.Name,
			Description: logItem.Description,
			Priority:    logItem.Priority,
			Removed:     logItem.Removed,
		})
	}

	if len(logs) == s.config.StackLogs {
		if err := s.log.Send(logs); err != nil {
			log.Warn(err)
		}

		logs = make([]*domain.ChangesHistory, 0, s.config.StackLogs)
	}
}

func (s *GoodsService) Shutdown() {
	if err := s.log.Send(logs); err != nil {
		log.Warn(err)
	}
}
