package user_rating_service

import (
	"context"
	"errors"

	"github.com/sosivvodnic/kinotower-go/internal/features/user_ratings/domain"
	user_rating_repository "github.com/sosivvodnic/kinotower-go/internal/features/user_ratings/repository"
)

var ErrValidation = errors.New("validation error")
var ErrUserNotFound = errors.New("user not found")
var ErrFilmNotFound = errors.New("film not found")
var ErrScoreExists = errors.New("score exists")
var ErrRatingNotFound = errors.New("rating not found")

type Service interface {
	Create(ctx context.Context, userID int, req domain.CreateRequest) (*domain.RatingResponse, error)
	List(ctx context.Context, userID int) ([]domain.RatingResponse, error)
	Delete(ctx context.Context, userID int, ratingID int) error
}

type service struct {
	repo user_rating_repository.Repository
}

func NewService(repo user_rating_repository.Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID int, req domain.CreateRequest) (*domain.RatingResponse, error) {
	if req.FilmID <= 0 {
		return nil, ErrValidation
	}
	if req.Ball < 1 || req.Ball > 5 {
		return nil, ErrValidation
	}

	ok, err := s.repo.UserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrUserNotFound
	}
	ok, err = s.repo.FilmExists(ctx, req.FilmID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrFilmNotFound
	}

	exists, err := s.repo.RatingExists(ctx, userID, req.FilmID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrScoreExists
	}

	return s.repo.Create(ctx, userID, req.FilmID, req.Ball)
}

func (s *service) List(ctx context.Context, userID int) ([]domain.RatingResponse, error) {
	ok, err := s.repo.UserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrUserNotFound
	}
	return s.repo.ListByUser(ctx, userID)
}

func (s *service) Delete(ctx context.Context, userID int, ratingID int) error {
	ok, err := s.repo.UserExists(ctx, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUserNotFound
	}
	if ratingID <= 0 {
		return ErrRatingNotFound
	}
	deleted, err := s.repo.Delete(ctx, userID, ratingID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrRatingNotFound
	}
	return nil
}

