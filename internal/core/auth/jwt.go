package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	Secret     []byte
	TokenTTL   time.Duration
	revokedJTI map[string]time.Time
}

func NewManager(secret string, ttl time.Duration) *Manager {
	return &Manager{
		Secret:     []byte(secret),
		TokenTTL:   ttl,
		revokedJTI: map[string]time.Time{},
	}
}

type Claims struct {
	UserID int    `json:"user_id"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

func (m *Manager) Issue(userID int) (token string, jti string, exp time.Time, err error) {
	jti, err = randomHex(16)
	if err != nil {
		return "", "", time.Time{}, err
	}
	exp = time.Now().UTC().Add(m.TokenTTL)

	claims := Claims{
		UserID: userID,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString(m.Secret)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return s, jti, exp, nil
}

func (m *Manager) Parse(token string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return m.Secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}

	// Cleanup and check revocation
	m.gcRevoked()
	if until, ok := m.revokedJTI[claims.JTI]; ok && until.After(time.Now().UTC()) {
		return nil, errors.New("token revoked")
	}

	return claims, nil
}

func (m *Manager) Revoke(jti string, until time.Time) {
	m.revokedJTI[jti] = until
	m.gcRevoked()
}

func (m *Manager) gcRevoked() {
	now := time.Now().UTC()
	for jti, until := range m.revokedJTI {
		if !until.After(now) {
			delete(m.revokedJTI, jti)
		}
	}
}

func randomHex(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
