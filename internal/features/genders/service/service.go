package gender_service

import (
	"context"

	"github.com/sosivvodnic/kinotower-go/internal/features/genders/domain"
	gender_repository "github.com/sosivvodnic/kinotower-go/internal/features/genders/repository"
)

type GenderService interface {
	List(ctx context.Context) ([]domain.Gender, error)
}

type genderService struct {
	repo gender_repository.GenderRepository
}

func NewGenderService(repo gender_repository.GenderRepository) *genderService {
	return &genderService{repo: repo}
}

func (s *genderService) List(ctx context.Context) ([]domain.Gender, error) {
	return s.repo.List(ctx)
}
