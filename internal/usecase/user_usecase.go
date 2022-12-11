package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (s *Storage) CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error) {
	s.logger.Debug("Enter in usecase CreateUser()")
	id, err := s.userStore.Create(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create user: %w", err)
	}
	return id, nil
}

