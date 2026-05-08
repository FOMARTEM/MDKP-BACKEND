package provider

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

/*
	Функции для работы с логами
	1) Создание лога (При любом действии пользователя, при входе в систему, при выходе из системы, при изменении данных аккаунта, при создании задачи, при удалении задачи, при изменении статуса задачи, при назначении задачи, при создании правки, при удалении правки, при изменении статуса правки)
	2) Получение логов по id пользователя
	3) Получение логов по почте пользователя
	4) Получение логов по определённым датам
*/

// Создание лога (При любом действии пользователя, при входе в систему, при выходе из системы, при изменении данных аккаунта, при создании задачи, при удалении задачи, при изменении статуса задачи, при назначении задачи, при создании правки, при удалении правки, при изменении статуса правки)
func (p *Provider) CreateLog(userID int, action string) error {
	_, err := p.conn.Exec(
		`INSERT INTO public."Log" ("UserID", "Action", "Date") VALUES ($1, $2, NOW())`,
		userID, action,
	)
	return err
}

// Получение логов по id пользователя
func (p *Provider) GetLogsByUserID(userID int) ([]entities.Log, error) {
	rows, err := p.conn.Query(
		`SELECT id, "UserID", "Action", "Date" FROM public."Log" WHERE "UserID" = $1 ORDER BY id DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []entities.Log
	for rows.Next() {
		var log entities.Log
		if err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.DateCreated,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}

// Получение логов по почте пользователя
func (p *Provider) GetLogsByUserEmail(email string) ([]entities.Log, error) {
	rows, err := p.conn.Query(
		`SELECT l.id, l."UserID", l."Action", l."Date"
		 FROM public."Log" l
		 JOIN public."Employee" e ON e.id = l."UserID"
		 WHERE lower(e."Email") = lower($1)
		 ORDER BY l.id DESC`,
		email,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []entities.Log
	for rows.Next() {
		var log entities.Log
		if err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.DateCreated,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}

// Получение логов по определённым датам
func (p *Provider) GetLogsByDateRange(startDate, endDate string) ([]entities.Log, error) {
	rows, err := p.conn.Query(
		`SELECT id, "UserID", "Action", "Date" FROM public."Log" WHERE "Date" BETWEEN $1 AND $2 ORDER BY id DESC`,
		startDate, endDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []entities.Log
	for rows.Next() {
		var log entities.Log
		if err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.DateCreated,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}

// Получение всех логов (с пагинацией)
func (p *Provider) GetLogsAll(limit, offset int) ([]entities.Log, error) {
	rows, err := p.conn.Query(
		`SELECT id, "UserID", "Action", "Date"
		 FROM public."Log"
		 ORDER BY id DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []entities.Log
	for rows.Next() {
		var log entities.Log
		if err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.DateCreated,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}
