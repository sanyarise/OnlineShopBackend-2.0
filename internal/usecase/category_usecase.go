package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

// CreateCategory call database method and returns id of created category or error
func (s *Storage) CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	log.Println("Enter in usecase CreateCategory()")
	id, err := s.categoryStore.CreateCategory(ctx, category)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create category: %w", err)
	}
	return id, nil
}

// GetCategoryList call database method and returns chan with all models.Category or error
func (s *Storage) GetCategoryList(ctx context.Context) (chan models.Category, error) {
	chin, err := s.categoryStore.GetCategoryList(ctx)
	if err != nil {
		return nil, err
	}
	chout := make(chan models.Category, 100)
	go func() {
		defer close(chout)
		for {
			select {
			case <-ctx.Done():
				return
			case category, ok := <-chin:
				if !ok {
					return
				}
				chout <- category
			}
		}
	}()
	return chout, nil

}
