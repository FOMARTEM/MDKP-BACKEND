package usecase

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

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
		// Без фильтров — отдаём все логи с пагинацией.
		return u.p.GetLogsAll(limit, offset)
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

func (u *Usecase) GetUsersByRole(user entities.User) ([]entities.User, error) {
	var searchBy string

	// Определяем по какому полю искать
	if user.Role != "" {
		searchBy = "position"
	} else if user.Email != "" {
		searchBy = "email"
	} else if user.LastName != "" {
		searchBy = "last_name"
	} else if user.FirstName != "" {
		searchBy = "first_name"
	} else {
		return nil, fmt.Errorf("no search criteria provided")
	}

	users, err := u.p.FindUsers(user, searchBy)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *Usecase) UserRoleUpdate(userEmail string, roleID int) error {
	user, err := u.p.GetUserByEmail(userEmail)
	if err != nil {
		return err
	}

	err = u.p.UpdateUserRole(user.ID, roleID)

	if err != nil {
		return err
	}

	return nil
}

func (u *Usecase) GetUsers() ([]entities.User, error) {
	users, err := u.p.ListUsers()

	if err != nil {
		return nil, err
	}

	return users, nil
}
