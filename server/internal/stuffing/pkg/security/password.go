package security

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type ISecurity interface {
	EncryptPassword(ctx context.Context, password string) (string, error)
	VerifyPassword(ctx context.Context, hashed, password string) error
}

type Security struct{}

func NewSecurity() ISecurity {
	return &Security{}
}

// Шифруем пас
func (sc *Security) EncryptPassword(ctx context.Context, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil

}

func (sc *Security) VerifyPassword(ctx context.Context, hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
