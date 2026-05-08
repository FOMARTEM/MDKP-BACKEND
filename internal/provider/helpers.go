package provider

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

// Получаем всевозможные роли пользователей системы из базы данных
func (p *Provider) GetRoles() ([]entities.Role, error) {
	rows, err := p.conn.Query(
		`SELECT id, "Name", "Description"
		 FROM public."Position"
		 ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []entities.Role
	for rows.Next() {
		var role entities.Role
		if err := rows.Scan(
			&role.ID,
			&role.Title,
			&role.Description,
		); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

// Получаем Id роли по её имени
func (p *Provider) GetRoleId(role string) (int, error) {
	var id int

	err := p.conn.QueryRow(
		`SELECT id
		 FROM public."Position"
		 where "Name" = $1`,
		role,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *Provider) GetStatuses() ([]entities.Status, error) {
	rows, err := p.conn.Query(
		`SELECT id, "Name", "Description"
		 FROM public."Status"
		 ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var statuses []entities.Status
	for rows.Next() {
		var status entities.Status
		if err := rows.Scan(
			&status.ID,
			&status.Title,
			&status.Description,
		); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	return statuses, rows.Err()
}
