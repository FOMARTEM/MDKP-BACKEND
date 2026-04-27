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

// Получение данных пользователя (без пароля)
func (u *Usecase) SelectUserByEmail(email string) (*entities.User, error) {
	user, err := u.p.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// Обновление пароля пользователя
func (u *Usecase) UpdateUserPassword(userID int, newPassword string) error {
	newPasswordHash, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	err = u.p.UpdateUserPasswordHash(userID, newPasswordHash)
	if err != nil {
		return err
	}

	err = u.p.CreateLog(userID, "Password updated")
	if err != nil {
		return err
	}

	return nil
}

// Получение данных пользователя (без пароля) по id
func (u *Usecase) SelectUserByID(userID int) (*entities.User, error) {
	user, err := u.p.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// Получение логов
func (u *Usecase) GetLogs(userID int, email string, startDate, endDate string, limit, offset int) ([]entities.Log, error) {
	var logs []entities.Log
	var err error

	if userID != 0 {
		logs, err = u.p.GetLogsByUserID(userID)
	} else if email != "" {
		logs, err = u.p.GetLogsByUserEmail(email)
	} else if startDate != "" && endDate != "" {
		logs, err = u.p.GetLogsByDateRange(startDate, endDate)
	} else {
		return nil, entities.ErrInvalidQueryParams
	}

	if err != nil {
		return nil, err
	}

	// Применяем пагинацию к результату
	if offset > len(logs) {
		return []entities.Log{}, nil
	}

	end := offset + limit
	if end > len(logs) {
		end = len(logs)
	}

	return logs[offset:end], nil
}

// Получение ролей
func (u *Usecase) GetRoles() ([]entities.Role, error) {
	var roles []entities.Role

	roles, err := u.p.GetRoles()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// Получение статусов
func (u *Usecase) GetStatuses() ([]entities.Status, error) {
	var statuses []entities.Status

	statuses, err := u.p.GetStatuses()
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

func (u *Usecase) CreateUser(user entities.User) (*entities.User, error) {
	var err error

	user.PasswordHash, err = hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.IsActive = true

	user.RoleID, err = u.p.GetRoleId(user.Role)

	if err != nil {
		return nil, err
	}

	var createdUser *entities.User
	createdUser, err = u.p.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil

}

func (u *Usecase) UserActiveChange(userEmail string) error {
	user, err := u.p.GetUserByEmail(userEmail)
	if err != nil {
		return err
	}

	err = u.p.ChangeUserActive(user.ID, !user.IsActive)
	if err != nil {
		return err
	}

	return nil
}
