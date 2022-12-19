package handlers

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ ICategoryHandlers = &CategoryHandlers{}

type CategoryHandlers struct {
	usecase usecase.ICategoryUsecase
	logger  *zap.Logger
}

func NewCategoryHandlers(usecase usecase.ICategoryUsecase, logger *zap.Logger) *CategoryHandlers {
	return &CategoryHandlers{usecase: usecase, logger: logger}
}

// Category is struct for DTO
type Category struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
}

// CreateCategory transform Category to models.Category and call usecase CreateCategory
func (handlers *CategoryHandlers) CreateCategory(ctx context.Context, category Category) (uuid.UUID, error) {
	handlers.logger.Debug("Enter in handlers CreateCategory()")
	newCategory := &models.Category{
		Name:        category.Name,
		Description: category.Description,
		Image:       category.Image,
	}
	id, err := handlers.usecase.CreateCategory(ctx, newCategory)
	if err != nil {
		return id, err
	}
	return id, nil
}

// UpdateCategory transform Category to models.Category and call usecase UpdateCategory
func (handlers *CategoryHandlers) UpdateCategory(ctx context.Context, category Category) error {
	handlers.logger.Debug("Enter in handlers UpdateCategory()")
	id, err := uuid.Parse(category.Id)
	if err != nil {
		return fmt.Errorf("invalid category uuid: %w", err)
	}

	updateCategory := &models.Category{
		Id:          id,
		Name:        category.Name,
		Description: category.Description,
		Image:       category.Image,
	}
	return handlers.usecase.UpdateCategory(ctx, updateCategory)
}

// GetCategory returns Category on id
func (handlers *CategoryHandlers) GetCategory(ctx context.Context, id string) (Category, error) {
	handlers.logger.Debug("Enter in handlers GetCategory()")
	uid, err := uuid.Parse(id)
	if err != nil {
		return Category{}, fmt.Errorf("invalid category uuid: %w", err)
	}
	category, err := handlers.usecase.GetCategory(ctx, uid)
	if err != nil {
		return Category{}, err
	}

	return Category{
		Id:          category.Id.String(),
		Name:        category.Name,
		Description: category.Description,
		Image:       category.Image,
	}, nil
}

// GetCategoryList returns list of all categories
func (handlers *CategoryHandlers) GetCategoryList(ctx context.Context) ([]Category, error) {
	handlers.logger.Debug("Enter in handlers GetCategoryList()")
	res := make([]Category, 0, 100)
	categories, err := handlers.usecase.GetCategoryList(ctx)
	if err != nil {
		return res, err
	}
	for _, category := range categories {
		res = append(res, Category{
			Id:          category.Id.String(),
			Name:        category.Name,
			Description: category.Description,
			Image:       category.Image,
		})
	}
	return res, nil
}

// DeleteCategory calls usecase method for delete category by id
func (handlers *CategoryHandlers) DeleteCategory(ctx context.Context, id uuid.UUID) (string, error) {
	handlers.logger.Debug("Enter in handlers DeleteCategory()")
	deletedCategoryName, err := handlers.usecase.DeleteCategory(ctx, id)
	if err != nil {
		return "", err
	}
	return deletedCategoryName, nil
}

func (handlers *CategoryHandlers) GetCategoryByName(ctx context.Context, name string) (Category, error) {
	handlers.logger.Debug("Enter in handlers GetCategoryByName()")
	category, err := handlers.usecase.GetCategoryByName(ctx, name)
	if err != nil {
		return Category{}, err
	}
	return Category{
		Id:          category.Id.String(),
		Name:        category.Name,
		Description: category.Description,
		Image:       category.Image,
	}, nil
}
