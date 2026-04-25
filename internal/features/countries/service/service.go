package country_service

import (
	"context"

	"github.com/sosivvodnic/kinotower-go/internal/features/countries/domain"
	country_repository "github.com/sosivvodnic/kinotower-go/internal/features/countries/repository"
)

type CountryService interface {
	List(ctx context.Context) ([]domain.Country, error)
}

type countryService struct {
	repo country_repository.CountryRepository
}

func NewCountryService(repo country_repository.CountryRepository) *countryService {
	return &countryService{repo: repo}
}

func (s *countryService) List(ctx context.Context) ([]domain.Country, error) {
	return s.repo.List(ctx)
}

