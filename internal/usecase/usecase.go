package usecase

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	p Provider
}

func NewUsecase(p Provider) *Usecase {
	return &Usecase{
		p: p,
	}
}

const maxPasswordHashLen = 511

// HashPassword хэширует пароль через bcrypt.
// Результат гарантированно помещается в 511 символов (bcrypt-хеш обычно 60 символов).
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedPassword)
	if len(hash) > maxPasswordHashLen {
		return "", entities.ErrPasswordTooLong
	}

	return hash, nil
}
