package user_service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/sosivvodnic/kinotower-go/internal/features/users/domain"
	user_repository "github.com/sosivvodnic/kinotower-go/internal/features/users/repository"
)

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

var ErrValidation = errors.New("validation error")

type Service interface {
	GetUserByID(ctx context.Context, id int) (*domain.UserResponse, error)
	UpdateUser(ctx context.Context, id int, req domain.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id int) error
}

type service struct {
	repo user_repository.Repository
}

func NewService(repo user_repository.Repository) *service {
	return &service{repo: repo}
}

func (s *service) GetUserByID(ctx context.Context, id int) (*domain.UserResponse, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *service) UpdateUser(ctx context.Context, id int, req domain.UpdateUserRequest) error {
	fio := strings.TrimSpace(req.FIO)
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if len(fio) < 2 || len(fio) > 150 {
		return ErrValidation
	}
	if len(email) < 4 || len(email) > 50 || !emailRe.MatchString(email) {
		return ErrValidation
	}
	if req.Birthday == "" {
		return ErrValidation
	}
	if req.GenderID <= 0 {
		return ErrValidation
	}

	// birthday string is validated at handler by time.Parse
	b := req.Birthday
	return s.repo.UpdateUser(ctx, id, fio, email, &b, req.GenderID)
}

func (s *service) DeleteUser(ctx context.Context, id int) error {
	return s.repo.SoftDeleteUser(ctx, id)
}
