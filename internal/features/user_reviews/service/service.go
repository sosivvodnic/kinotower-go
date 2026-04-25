package user_review_service

import (
	"context"
	"errors"
	"strings"

	"github.com/sosivvodnic/kinotower-go/internal/features/user_reviews/domain"
	user_review_repository "github.com/sosivvodnic/kinotower-go/internal/features/user_reviews/repository"
)

var ErrValidation = errors.New("validation error")
var ErrUserNotFound = errors.New("user not found")
var ErrFilmNotFound = errors.New("film not found")
var ErrReviewNotFound = errors.New("review not found")

type Service interface {
	Create(ctx context.Context, userID int, req domain.CreateRequest) (*domain.ReviewResponse, error)
	List(ctx context.Context, userID int) ([]domain.ReviewResponse, error)
	Delete(ctx context.Context, userID int, reviewID int) error
}

type service struct {
	repo user_review_repository.Repository
}

func NewService(repo user_review_repository.Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID int, req domain.CreateRequest) (*domain.ReviewResponse, error) {
	if req.FilmID <= 0 {
		return nil, ErrValidation
	}
	msg := strings.TrimSpace(req.Message)
	if len(msg) < 4 || len(msg) > 1024 {
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

	return s.repo.Create(ctx, userID, req.FilmID, msg)
}

func (s *service) List(ctx context.Context, userID int) ([]domain.ReviewResponse, error) {
	ok, err := s.repo.UserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrUserNotFound
	}
	return s.repo.ListByUser(ctx, userID)
}

func (s *service) Delete(ctx context.Context, userID int, reviewID int) error {
	ok, err := s.repo.UserExists(ctx, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUserNotFound
	}
	if reviewID <= 0 {
		return ErrReviewNotFound
	}

	deleted, err := s.repo.Delete(ctx, userID, reviewID)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrReviewNotFound
	}
	return nil
}
