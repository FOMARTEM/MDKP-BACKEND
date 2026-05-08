package provider

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

/*
	Фукнции для работы с правками
	1) Создание правки
	2) Получение правок по id задачи
	3) Получение правки по id правки
	4) Изменение статуса правки (Открыта, В работе, На проверке, Закрыта)
*/

// Создание правки
func (p *Provider) CreateRevision(revision entities.Revision) (*entities.Revision, error) {
	var id int
	err := p.conn.QueryRow(
		`INSERT INTO public."Revision" ("CreateDate","Title","Description","EmployeeID","VersionID","StatusID")
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		revision.DateCreated,
		revision.Title,
		revision.Description,
		revision.CreatorID,
		revision.VersionID,
		revision.StatusID,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &entities.Revision{
		ID:          id,
		DateCreated: revision.DateCreated,
		Title:       revision.Title,
		Description: revision.Description,
		CreatorID:   revision.CreatorID,
		VersionID:   revision.VersionID,
		StatusID:    revision.StatusID,
	}, nil
}

// Получение правок по id задачи
func (p *Provider) GetRevisionsByVersionID(versionID int) ([]entities.Revision, error) {
	rows, err := p.conn.Query(
		`SELECT r.id, r."CreateDate", r."Title", r."Description", r."EmployeeID", r."VersionID", r."StatusID"
		 FROM public."Revision" r
		 WHERE r."VersionID" = $1
		 ORDER BY r.id`,
		versionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var revisions []entities.Revision
	for rows.Next() {
		var revision entities.Revision
		if err := rows.Scan(
			&revision.ID,
			&revision.DateCreated,
			&revision.Title,
			&revision.Description,
			&revision.CreatorID,
			&revision.VersionID,
			&revision.StatusID,
		); err != nil {
			return nil, err
		}
		revisions = append(revisions, revision)
	}
	return revisions, rows.Err()
}

// Получение правки по id правки
func (p *Provider) GetRevisionByID(revisionID int) (*entities.Revision, error) {
	var revision entities.Revision
	err := p.conn.QueryRow(
		`SELECT id, "CreateDate", "Title", "Description", "EmployeeID", "VersionID", "StatusID"
		 FROM public."Revision"	
	 WHERE id = $1`,
		revisionID,
	).Scan(
		&revision.ID,
		&revision.DateCreated,
		&revision.Title,
		&revision.Description,
		&revision.CreatorID,
		&revision.VersionID,
		&revision.StatusID,
	)

	if err != nil {
		return nil, err
	}

	return &revision, nil
}

// Изменение статуса правки (Открыта, В работе, На проверке, Закрыта)
func (p *Provider) UpdateRevisionStatus(revisionID int, status string) error {
	_, err := p.conn.Exec(
		`UPDATE public."Revision" SET "StatusID" = (SELECT id FROM public."Status" WHERE "Name" = $2) WHERE id = $1`,
		revisionID, status,
	)
	return err
}
