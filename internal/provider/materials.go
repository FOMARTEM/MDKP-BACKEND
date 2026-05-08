package provider

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

/*
	Функции для работы с материалами
	1) Создание материала
	2) Получение материалов по id задачи
	3) Получение материала по id материала
	4) Присвоение материала к задаче
	5) Присвоение материала к версии
*/

// Создание материала
func (p *Provider) CreateMaterial(material entities.Material) (*entities.Material, error) {
	var id int
	err := p.conn.QueryRow(
		`INSERT INTO public."Material" ("FileName","Extension","Description","EmployeeID","TaskID")
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		material.Title,
		material.Extension,
		material.Description,
		material.CreatorID,
		material.TaskID,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &entities.Material{
		ID:          id,
		Title:       material.Title,
		Extension:   material.Extension,
		Description: material.Description,
		CreatorID:   material.CreatorID,
		TaskID:      material.TaskID,
	}, nil
}

// Получение материалов по id задачи
func (p *Provider) GetMaterialsByTaskID(taskID int) ([]entities.Material, error) {
	rows, err := p.conn.Query(
		`SELECT id, "FileName", "Extension", "Description", "EmployeeID", "TaskID"
		 FROM public."Material"
		 WHERE "TaskID" = $1
		 ORDER BY id`,
		taskID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var materials []entities.Material
	for rows.Next() {
		var material entities.Material
		if err := rows.Scan(
			&material.ID,
			&material.Title,
			&material.Extension,
			&material.Description,
			&material.CreatorID,
			&material.TaskID,
		); err != nil {
			return nil, err
		}
		materials = append(materials, material)
	}
	return materials, rows.Err()
}

// Получение материала по id материала
func (p *Provider) GetMaterialByID(materialID int) (*entities.Material, error) {
	var material entities.Material
	err := p.conn.QueryRow(
		`SELECT id, "FileName", "Extension", "Description", "EmployeeID", "TaskID"
		 FROM public."Material"
		 WHERE id = $1`,
		materialID,
	).Scan(
		&material.ID,
		&material.Title,
		&material.Extension,
		&material.Description,
		&material.CreatorID,
		&material.TaskID,
	)

	if err != nil {
		return nil, err
	}

	return &material, nil
}

// Присвоение материала к задаче
func (p *Provider) AssignMaterialToTask(materialID int, taskID int) error {
	_, err := p.conn.Exec(
		`UPDATE public."Material" SET "TaskID" = $2 WHERE id = $1`,
		materialID, taskID,
	)
	return err
}

// Присвоение материала к версии
func (p *Provider) AssignMaterialToVersion(materialID int, versionID int) error {
	_, err := p.conn.Exec(
		`UPDATE public."Material" SET "VersionID" = $2 WHERE id = $1`,
		materialID, versionID,
	)
	return err
}
