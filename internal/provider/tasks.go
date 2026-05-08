package provider

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	"github.com/lib/pq"
)

/*
	Функции для работы с задачами
	1) Создание задачи
	2) Получение списка задач пользователя (по id или почте)
	3) Получение задачи по id
	4) Изменение статуса задачи (Открыта, В работе, На проверке, Закрыта)
	5) Изменение автора задачи
	6) Изменение редактора задачи
	7) Установка даты готовности задачи
	8) Получение задач по приоритетам
	9) Получение задач с определёнными статусами
*/

// Создание задачи
func (p *Provider) CreateTask(task entities.Tasks) (*entities.Tasks, error) {
	var id int
	err := p.conn.QueryRow(
		`INSERT INTO public."Task" ("Title","Description","CreateDate","DeadlineDate","Priority","CreatorId","EditorID","AuthorID","StatusID")
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING "id"`,
		task.Title,
		task.Description,
		task.DateCreated,
		task.DateDedline,
		task.Priority,
		task.IdCreator,
		task.IdRedactor,
		task.IdAuthor,
		task.IdStatus,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &entities.Tasks{
		ID:          id,
		Title:       task.Title,
		Description: task.Description,
		DateCreated: task.DateCreated,
		DateDedline: task.DateDedline,
		DateClosed:  task.DateClosed,
		Priority:    task.Priority,
		IdCreator:   task.IdCreator,
		IdRedactor:  task.IdRedactor,
		IdAuthor:    task.IdAuthor,
		IdStatus:    task.IdStatus,
	}, nil
}

func (p *Provider) DeleteTask(id int) error {
	_, err := p.conn.Exec(
		`DELETE FROM public."Task" WHERE id = $1`,
		id,
	)
	return err
}

// Получение списка задач пользователя (по id или почте)
func (p *Provider) GetTasksByUserID(userID int) ([]entities.Tasks, error) {
	rows, err := p.conn.Query(
		`SELECT 			
			id, 
			"Title", 
			"Description", 
			"CreateDate", 
			"DeadlineDate", 
			COALESCE("ReadyDate"::text, '') as "ReadyDate", 
			"Priority", 
			"CreatorId", 
			"EditorID", 
			"AuthorID", 
			"StatusID"
		FROM public."Task" t
		WHERE t."CreatorId" = $1 
			OR t."EditorID" = $1 
			OR t."AuthorID" = $1
		ORDER BY t.id`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []entities.Tasks
	for rows.Next() {
		var task entities.Tasks
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DateCreated,
			&task.DateDedline,
			&task.DateClosed,
			&task.Priority,
			&task.IdCreator,
			&task.IdRedactor,
			&task.IdAuthor,
			&task.IdStatus,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// Получение задачи по id
func (p *Provider) GetTaskByID(taskID int) (*entities.Tasks, error) {
	var task entities.Tasks
	err := p.conn.QueryRow(
		`SELECT 
			id, 
			"Title", 
			"Description", 
			"CreateDate", 
			"DeadlineDate", 
			COALESCE("ReadyDate"::text, '') as "ReadyDate", 
			"Priority", 
			"CreatorId", 
			"EditorID", 
			"AuthorID", 
			"StatusID"
		FROM public."Task"
		WHERE id = $1`,
		taskID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DateCreated,
		&task.DateDedline,
		&task.DateClosed,
		&task.Priority,
		&task.IdCreator,
		&task.IdRedactor,
		&task.IdAuthor,
		&task.IdStatus,
	)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Изменение статуса задачи (Открыта, В работе, На проверке, Закрыта)
func (p *Provider) UpdateTaskStatus(taskID int, status string) error {
	_, err := p.conn.Exec(
		`UPDATE public."Task" SET "StatusID" = (SELECT id FROM public."Status" WHERE "Name" = $2) WHERE id = $1`,
		taskID, status,
	)
	return err
}

// Изменение автора задачи
func (p *Provider) UpdateTaskAuthor(taskID int, authorID int, authorEmail string) error {
	_, err := p.conn.Exec(
		`UPDATE public."Task" SET "AuthorID" = (SELECT id FROM public."Employee" WHERE id = $2 OR lower("Email") = lower($3)) WHERE id = $1`,
		taskID, authorID, authorEmail,
	)
	return err
}

// Изменение редактора задачи
func (p *Provider) UpdateTaskEditor(taskID int, editorID int, editorEmail string) error {
	_, err := p.conn.Exec(
		`UPDATE public."Task" SET "EditorID" = (SELECT id FROM public."Employee" WHERE id = $2 OR lower("Email") = lower($3)) WHERE id = $1`,
		taskID, editorID, editorEmail,
	)
	return err
}

// Установка даты готовности задачи
func (p *Provider) UpdateTaskReadyDate(taskID int, readyDate string) error {
	_, err := p.conn.Exec(
		`UPDATE public."Task" SET "ReadyDate" = $2 WHERE id = $1`,
		taskID, readyDate,
	)
	return err
}

// Получение задач по приоритетам
func (p *Provider) GetTasksByPriority(priority int) ([]entities.Tasks, error) {
	rows, err := p.conn.Query(
		`SELECT 
			id, 
			"Title", 
			"Description", 
			"CreateDate", 
			"DeadlineDate", 
			COALESCE("ReadyDate"::text, '') as "ReadyDate", 
			"Priority", 
			"CreatorId", 
			"EditorID", 
			"AuthorID", 
			"StatusID"
		 FROM public."Task"
		 WHERE "Priority" = $1
		 ORDER BY id`,
		priority,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []entities.Tasks
	for rows.Next() {
		var task entities.Tasks
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DateCreated,
			&task.DateDedline,
			&task.DateClosed,
			&task.Priority,
			&task.IdCreator,
			&task.IdRedactor,
			&task.IdAuthor,
			&task.IdStatus,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// Получение задач с определёнными статусами
func (p *Provider) GetTasksByStatuses(statuses []string) ([]entities.Tasks, error) {
	if len(statuses) == 0 {
		return []entities.Tasks{}, nil
	}
	query := `
		SELECT 
			id, 
			"Title", 
			"Description", 
			"CreateDate", 
			"DeadlineDate", 
			COALESCE("ReadyDate"::text, '') as "ReadyDate", 
			"Priority", 
			"CreatorId", 
			"EditorID", 
			"AuthorID", 
			"StatusID"
		FROM public."Task"
		WHERE "StatusID" IN (SELECT id FROM public."Status" WHERE "Name" = ANY($1))
		ORDER BY id`
	rows, err := p.conn.Query(query, pq.Array(statuses))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []entities.Tasks
	for rows.Next() {
		var task entities.Tasks
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DateCreated,
			&task.DateDedline,
			&task.DateClosed,
			&task.Priority,
			&task.IdCreator,
			&task.IdRedactor,
			&task.IdAuthor,
			&task.IdStatus,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}
