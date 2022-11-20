package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Storage) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	id, err := s.store.CreateItem(ctx, item)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create item: %w", err)
	}
	return id, nil
}

func (s *Storage) UpdateItem(ctx context.Context, item *models.Item) error {
	return s.store.UpdateItem(ctx, item)
}

func (s *Storage) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	item, err := s.store.GetItem(ctx, id)
	if err != nil {
		return &models.Item{}, fmt.Errorf("error on get item: %w", err)
	}
	return item, nil
}

func (s *Storage) ItemsList(ctx context.Context) (chan models.Item, error) {
	chin, err := s.store.ItemsList(ctx)
	if err != nil {
		return nil, err
	}
	chout := make(chan models.Item, 100)
	go func() {
		defer close(chout)
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-chin:
				if !ok {
					return
				}
				chout <- item
			}
		}
	}()
	return chout, nil

}

func (s *Storage) SearchLine(ctx context.Context, param string) (chan models.Item, error) {
	chin, err := s.store.SearchLine(ctx, param)
	if err != nil {
		return nil, err
	}
	chout := make(chan models.Item, 100)
	go func() {
		defer close(chout)
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-chin:
				if !ok {
					return
				}
				chout <- item
			}
		}
	}()
	return chout, nil
}
