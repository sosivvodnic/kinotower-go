package category_service

import (
	"context"

	"github.com/sosivvodnic/kinotower-go/internal/features/categories/domain"
	category_repository "github.com/sosivvodnic/kinotower-go/internal/features/categories/repository"
)

type CategoryService interface {
	List(ctx context.Context) ([]domain.Category, error)
}

type categoryService struct {
	repo category_repository.CategoryRepository
}

func NewCategoryService(repo category_repository.CategoryRepository) *categoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) List(ctx context.Context) ([]domain.Category, error) {
	return s.repo.List(ctx)
}
