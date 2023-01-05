package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository"
	"OnlineShopBackend/internal/repository/cash"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ ICategoryUsecase = &CategoryUsecase{}

var (
	categoriesListKey = "CategoriesList"
)

type CategoryUsecase struct {
	categoryStore  repository.CategoryStore
	categoriesCash cash.ICategoriesCash
	logger         *zap.Logger
}

func NewCategoryUsecase(store repository.CategoryStore, cash cash.ICategoriesCash, logger *zap.Logger) ICategoryUsecase {
	logger.Debug("Enter in usecase NewCategoryUsecase()")
	return &CategoryUsecase{categoryStore: store, categoriesCash: cash, logger: logger}
}

// / CreateCategory call database method and returns id of created category or error
func (usecase *CategoryUsecase) CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase CreateCategory() with args: ctx, category: %v", category)
	id, err := usecase.categoryStore.CreateCategory(ctx, category)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create category: %w", err)
	}
	err = usecase.UpdateCash(ctx, id, "create")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cash: %v", err))
	} else {
		usecase.logger.Info("Update cash success")
	}
	return id, nil
}

// UpdateCategory call database method to update category and returns error or nil
func (usecase *CategoryUsecase) UpdateCategory(ctx context.Context, category *models.Category) error {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateCategory() with args: ctx, category: %v", category)
	err := usecase.categoryStore.UpdateCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("error on update category: %w", err)
	}
	err = usecase.UpdateCash(ctx, category.Id, "update")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cash: %v", err))
	} else {
		usecase.logger.Info("Update cash success")
	}
	return nil
}

// GetCategory call database and returns *models.Category with given id or returns error
func (usecase *CategoryUsecase) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetCategory() with args: ctx, id: %v", id)
	category, err := usecase.categoryStore.GetCategory(ctx, id)
	if err != nil {
		return &models.Category{}, fmt.Errorf("error on get category: %w", err)
	}
	return category, nil
}

// GetCategoryList call database method and returns chan with all models.Category or error
func (usecase *CategoryUsecase) GetCategoryList(ctx context.Context) ([]models.Category, error) {
	usecase.logger.Debug("Enter in usecase GetCategoryList() with args: ctx")
	if ok := usecase.categoriesCash.CheckCash(ctx, categoriesListKey); !ok {
		categoryIncomingChan, err := usecase.categoryStore.GetCategoryList(ctx)
		if err != nil {
			return nil, err
		}
		categories := make([]models.Category, 0, 100)
		for category := range categoryIncomingChan {
			usecase.logger.Debug(fmt.Sprintf("category from channel is: %v", category))
			categories = append(categories, category)
		}
		err = usecase.categoriesCash.CreateCategoriesListCash(ctx, categories, categoriesListKey)
		if err != nil {
			return nil, fmt.Errorf("error on create categories list cash: %w", err)
		}
	}

	categories, err := usecase.categoriesCash.GetCategoriesListCash(ctx, categoriesListKey)
	if err != nil {
		return nil, fmt.Errorf("error on get cash: %w", err)
	}
	usecase.logger.Info("Get category list success")
	return categories, nil
}

// DeleteCategory call database method for deleting category
func (usecase *CategoryUsecase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteCategory() with args: ctx, id: %v", id)
	err := usecase.categoryStore.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}
	err = usecase.UpdateCash(ctx, id, "delete")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cash: %v", err))
	}
	usecase.logger.Info("Delete category success")
	return nil
}

func (usecase *CategoryUsecase) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetCategoryByName() with args: ctx, name: %s", name)
	category, err := usecase.categoryStore.GetCategoryByName(ctx, name)
	if err != nil {
		return nil, err
	}
	usecase.logger.Info("Get category by name success")
	return category, nil
}

// UpdateCash updating cash when creating or updating category
func (usecase *CategoryUsecase) UpdateCash(ctx context.Context, id uuid.UUID, op string) error {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateCash() with args: ctx, id: %v, op: %s", id, op)
	if !usecase.categoriesCash.CheckCash(ctx, categoriesListKey) {
		return fmt.Errorf("cash is not exists")
	}
	newCategory, err := usecase.categoryStore.GetCategory(ctx, id)
	if err != nil {
		if op == "delete" {
			newCategory = &models.Category{Id: id}
		} else {
			return fmt.Errorf("error on get category: %w", err)
		}
	}
	categories, err := usecase.categoriesCash.GetCategoriesListCash(ctx, categoriesListKey)
	if err != nil {
		return fmt.Errorf("error on get cash: %w", err)
	}
	if op == "update" {
		for i, category := range categories {
			if category.Id == id {
				categories[i] = *newCategory
				break
			}
		}
	}
	if op == "create" {
		categories = append(categories, *newCategory)
	}
	if op == "delete" {
		for i, category := range categories {
			if category.Id == id {
				categories = append(categories[:i], categories[i+1:]...)
				break
			}
		}
	}
	err = usecase.categoriesCash.CreateCategoriesListCash(ctx, categories, categoriesListKey)
	if err != nil {
		return err
	}
	usecase.logger.Info("Category cash update success")
	return nil
}

// DeleteCategoryCash deleted cash by deleting categories
func (usecase *CategoryUsecase) DeleteCategoryCash(ctx context.Context, name string) error {
	usecase.logger.Debug(fmt.Sprintf("Enter in usecase DeleteCategoryCash() with args: ctx, name: %s", name))
	err := usecase.categoriesCash.DeleteCash(ctx, name)
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on delete cash with key: %s, error is %v", name, err))
		return err
	}
	err = usecase.categoriesCash.DeleteCash(ctx, name+"Quantity")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on delete cash with key: %s, error is %v", name, err))
		return err
	}
	usecase.logger.Info("Category cash deleted success")
	return nil
}
