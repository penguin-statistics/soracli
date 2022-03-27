package services

import (
	"context"
	"time"

	"github.com/penguin-statistics/soracli/internal/models"
	"github.com/penguin-statistics/soracli/internal/models/cache"
	"github.com/penguin-statistics/soracli/internal/models/types"
	"github.com/penguin-statistics/soracli/internal/pkg/client"
)

type ItemService struct {
	http *client.HTTP
}

func NewItemService(http *client.HTTP) *ItemService {
	return &ItemService{
		http: http,
	}
}

func (s *ItemService) GetItems(ctx context.Context) ([]*models.Item, error) {
	var resp types.CliGameDataSeedResponse
	err := cache.CliGameDataSeed.Get(&resp)
	if err == nil {
		return resp.Items, nil
	}

	err = s.http.GetJSON("/cli/gamedata/seed", &resp)
	if err != nil {
		return nil, err
	}
	cache.CliGameDataSeed.Set(resp, 24*time.Hour)
	return resp.Items, nil
}

func (s *ItemService) GetItemsMapByArkId(ctx context.Context) (map[string]*models.Item, error) {
	var itemsMapByArkId map[string]*models.Item
	err := cache.ItemsMapByArkID.MutexGetSet(&itemsMapByArkId, func() (map[string]*models.Item, error) {
		items, err := s.GetItems(ctx)
		if err != nil {
			return nil, err
		}
		s := make(map[string]*models.Item)
		for _, item := range items {
			s[item.ArkItemID] = item
		}
		return s, nil
	}, 24*time.Hour)
	if err != nil {
		return nil, err
	}
	return itemsMapByArkId, nil
}
