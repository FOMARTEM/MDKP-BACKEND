package usecase

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

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
