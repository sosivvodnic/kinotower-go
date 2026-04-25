package auth_service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/sosivvodnic/kinotower-go/internal/core/auth"
	auth_repository "github.com/sosivvodnic/kinotower-go/internal/features/auth/repository"
	"golang.org/x/crypto/bcrypt"
)

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrValidation = errors.New("validation error")

type Service interface {
	SignUp(ctx context.Context, fio, email, password string, birthday *time.Time, genderID int) (token string, id int, outFIO string, err error)
	SignIn(ctx context.Context, email, password string) (token string, id int, fio string, jti string, exp time.Time, err error)
}

type service struct {
	repo auth_repository.Repository
	jwt  *auth.Manager
}

func NewService(repo auth_repository.Repository, jwt *auth.Manager) *service {
	return &service{repo: repo, jwt: jwt}
}

func (s *service) SignUp(ctx context.Context, fio, email, password string, birthday *time.Time, genderID int) (string, int, string, error) {
	fio = strings.TrimSpace(fio)
	email = strings.TrimSpace(strings.ToLower(email))

	if len(fio) < 2 || len(fio) > 150 {
		return "", 0, "", ErrValidation
	}
	if len(email) < 4 || len(email) > 50 || !emailRe.MatchString(email) {
		return "", 0, "", ErrValidation
	}
	if len(password) < 6 || len(password) > 65536 {
		return "", 0, "", ErrValidation
	}
	if genderID <= 0 {
		return "", 0, "", ErrValidation
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", 0, "", err
	}

	id, outFIO, err := s.repo.CreateUser(ctx, fio, email, string(hash), birthday, genderID)
	if err != nil {
		return "", 0, "", err
	}

	token, _, _, err := s.jwt.Issue(id)
	if err != nil {
		return "", 0, "", err
	}
	return token, id, outFIO, nil
}

func (s *service) SignIn(ctx context.Context, email, password string) (string, int, string, string, time.Time, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return "", 0, "", "", time.Time{}, ErrInvalidCredentials
	}

	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", 0, "", "", time.Time{}, err
	}
	if u == nil {
		return "", 0, "", "", time.Time{}, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", 0, "", "", time.Time{}, ErrInvalidCredentials
	}

	token, jti, exp, err := s.jwt.Issue(u.ID)
	if err != nil {
		return "", 0, "", "", time.Time{}, err
	}
	return token, u.ID, u.FIO, jti, exp, nil
}

