package provider

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

/*
	Функции для работы с версиями
	1) Создание версии задачи
	2) Получение версий задачи по id задачи
	3) Получение версии задачи по id версии
	4) Изменение статуса версии (Открыта, В работе, На проверке, Закрыта)
*/

// Создание версии задачи
func (p *Provider) CreateVersion(version entities.Version) (*entities.Version, error) {
	var id int
	err := p.conn.QueryRow(
		`INSERT INTO public."Version" ("VersionNumber","CreateDate","Title","Description","EmployeeID","TaskID")
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		version.NumberVersion,
		version.DateCreated,
		version.Title,
		version.Description,
		version.CreatorID,
		version.TaskID,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &entities.Version{
		ID:            id,
		NumberVersion: version.NumberVersion,
		DateCreated:   version.DateCreated,
		Title:         version.Title,
		Description:   version.Description,
		CreatorID:     version.CreatorID,
		TaskID:        version.TaskID,
	}, nil
}

// Получение версий задачи по id задачи
func (p *Provider) GetVersionsByTaskID(taskID int) ([]entities.Version, error) {
	rows, err := p.conn.Query(
		`SELECT id, "VersionNumber", "CreateDate", "Title", "Description", "EmployeeID", "TaskID"
		 FROM public."Version"
		 WHERE "TaskID" = $1
		 ORDER BY id`,
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var versions []entities.Version
	for rows.Next() {
		var version entities.Version
		if err := rows.Scan(
			&version.ID,
			&version.NumberVersion,
			&version.DateCreated,
			&version.Title,
			&version.Description,
			&version.CreatorID,
			&version.TaskID,
		); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, rows.Err()
}

// Получение версии задачи по id версии
func (p *Provider) GetVersionByID(versionID int) (*entities.Version, error) {
	var version entities.Version
	err := p.conn.QueryRow(
		`SELECT id, "VersionNumber", "CreateDate", "Title", "Description", "EmployeeID", "TaskID"
		 FROM public."Version"
		 WHERE id = $1`,
		versionID,
	).Scan(
		&version.ID,
		&version.NumberVersion,
		&version.DateCreated,
		&version.Title,
		&version.Description,
		&version.CreatorID,
		&version.TaskID,
	)

	if err != nil {
		return nil, err
	}

	return &version, nil
}

// Изменение статуса версии (Открыта, В работе, На проверке, Закрыта)
func (p *Provider) UpdateVersionStatus(versionID int, status string) error {
	_, err := p.conn.Exec(
		`UPDATE public."Version" SET "StatusID" = (SELECT id FROM public."Status" WHERE "Name" = $2) WHERE id = $1`,
		versionID, status,
	)
	return err
}
