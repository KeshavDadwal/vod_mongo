package auth

import (
	"crypto/sha256"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/o1egl/paseto"
)

func GenerateSignupToken(userID, username, email string) (string, error) {
	secret := os.Getenv("PASETO_SECRET")
	if secret == "" {
		return "", errors.New("missing PASETO_SECRET")
	}

	ttl := 60 * time.Minute
	if ttlRaw := os.Getenv("PASETO_TTL_MINUTES"); ttlRaw != "" {
		if n, err := strconv.Atoi(ttlRaw); err == nil && n > 0 {
			ttl = time.Duration(n) * time.Minute
		}
	}

	now := time.Now().UTC()
	claims := paseto.JSONToken{
		Subject:    userID,
		IssuedAt:   now,
		Expiration: now.Add(ttl),
	}
	claims.Set("username", username)
	claims.Set("email", email)

	// v2.local requires a 32-byte symmetric key.
	key := sha256.Sum256([]byte(secret))
	return paseto.NewV2().Encrypt(key[:], claims, nil)
}

