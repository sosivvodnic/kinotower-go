package auth_handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sosivvodnic/kinotower-go/internal/core/auth"
	"github.com/sosivvodnic/kinotower-go/internal/core/httpjson"
	"github.com/sosivvodnic/kinotower-go/internal/features/auth/domain"
	auth_service "github.com/sosivvodnic/kinotower-go/internal/features/auth/service"
)

type Handler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	SignOut(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	svc auth_service.Service
	jwt *auth.Manager
}

func NewHandler(svc auth_service.Service, jwt *auth.Manager) *handler {
	return &handler{svc: svc, jwt: jwt}
}

func (h *handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req domain.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "Bad request")
		return
	}

	var birthday *time.Time
	if req.Birthday != "" {
		t, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			httpjson.WriteError(w, http.StatusBadRequest, "Bad request")
			return
		}
		birthday = &t
	}

	token, id, fio, err := h.svc.SignUp(r.Context(), req.FIO, req.Email, req.Password, birthday, req.GenderID)
	if err != nil {
		if err == auth_service.ErrValidation {
			httpjson.WriteError(w, http.StatusBadRequest, "Bad request")
			return
		}
		// unique email -> 500/400? spec doesn't define; return 400
		httpjson.WriteError(w, http.StatusBadRequest, "Bad request")
		return
	}

	httpjson.WriteJSON(w, http.StatusCreated, domain.AuthResponse{
		Status: "success",
		Token:  token,
		ID:     id,
		FIO:    fio,
	})
}

func (h *handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req domain.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpjson.WriteError(w, http.StatusBadRequest, "Bad request")
		return
	}

	token, id, fio, _, _, err := h.svc.SignIn(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == auth_service.ErrInvalidCredentials {
			httpjson.WriteJSON(w, http.StatusUnauthorized, domain.InvalidAuthResponse{
				Status:  "invalid",
				Message: "Wrong email or password",
			})
			return
		}
		httpjson.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	httpjson.WriteJSON(w, http.StatusOK, domain.AuthResponse{
		Status: "success",
		Token:  token,
		ID:     id,
		FIO:    fio,
	})
}

func (h *handler) SignOut(w http.ResponseWriter, r *http.Request) {
	// Parse token again to extract JTI/exp and revoke it in memory.
	authHeader := r.Header.Get("Authorization")
	claims, err := authHeaderClaims(h.jwt, authHeader)
	if err == nil {
		until := time.Now().UTC()
		if claims.ExpiresAt != nil {
			until = claims.ExpiresAt.Time
		}
		h.jwt.Revoke(claims.JTI, until)
	}
	httpjson.WriteJSON(w, http.StatusOK, map[string]any{"status": "success"})
}

func authHeaderClaims(m *auth.Manager, h string) (*auth.Claims, error) {
	// expects "Bearer <token>"
	const pfx = "Bearer "
	if len(h) <= len(pfx) || h[:len(pfx)] != pfx {
		return nil, http.ErrNoCookie
	}
	return m.Parse(h[len(pfx):])
}

