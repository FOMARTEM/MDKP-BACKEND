package usecase

import (
	"database/sql"
	"errors"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	"golang.org/x/crypto/bcrypt"
)

// Функция для авторизации пользователя
func (u *Usecase) Authorize(user entities.User) (*bool, error) {
	login := false

	dbUser, err := u.p.GetUserByEmail(user.Email)

	if err != nil {
		return nil, entities.ErrUserNotFound
	}

	userPasswordHash, err := u.p.GetUserPasswordHashByEmail(user.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = u.p.CreateLog(dbUser.ID, "Failed login attempt - user not found")
			if err != nil {
				return nil, err
			}
			return &login, entities.ErrUserNotFound
		}
		return nil, err
	}

	if userPasswordHash != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(userPasswordHash), []byte(user.Password)); err != nil {
			err = u.p.CreateLog(dbUser.ID, "Failed login attempt")
			if err != nil {
				return nil, err
			}
			return &login, entities.ErrInvalidPassword
		}
		login = true
	} else {
		login = false
		err = u.p.CreateLog(dbUser.ID, "Failed login attempt - user not found")
		if err != nil {
			return nil, err
		}
		return &login, entities.ErrUserNotFound
	}

	err = u.p.CreateLog(dbUser.ID, "Successful login")
	if err != nil {
		return nil, err
	}

	return &login, nil
}
