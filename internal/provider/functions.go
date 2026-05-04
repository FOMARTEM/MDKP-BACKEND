package provider

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	"github.com/lib/pq"
)

/*
	Пользовательские функции для работы с базой данных
	1) Создаение пользователя
	2) Получение данных пользователя (без пароля)
	3) Получение хеша пароля по почте
	4) Получение хеша пароля по id
	5) Получение прав пользователя по id или почте
	6) Изменение хеша пароля пользователя
	7) Удаление пользователя (Изменяем IsActive на false)
	8) Поиск пользователя по почте (фио, почта, права, телефон)
	9) Поиск пользователя по правам (Администратор, Руководитель, Редактор, Автор)
	10) Изменение прав пользователя (Администратор, Руководитель, Редактор, Автор)
*/

// Создание пользователя
func (p *Provider) CreateUser(user entities.User) (*entities.User, error) {
	var id int

	err := p.conn.QueryRow(
		`INSERT INTO public."Employee"
		 ("LastName","FirstName","MiddleName","Email","Phone","BirthDate","PasswordHash","IsActive","Position_ID")
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING "id"`,
		user.LastName, user.FirstName, user.MiddleName, user.Email,
		user.Phone, user.DateOfBirth, user.PasswordHash, user.IsActive, user.RoleID,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:         id,
		LastName:   user.LastName,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		Email:      user.Email,
	}, nil
}

// Получение данных пользователя (без пароля)
func (p *Provider) GetUserByID(userID int) (*entities.User, error) {
	var user entities.User

	err := p.conn.QueryRow(
		`SELECT e.id, e."LastName", e."FirstName", e."MiddleName", e."Phone", e."BirthDate", e."Email", e."IsActive", p.id, p."Name"
		 FROM public."Employee" e
		 LEFT JOIN public."Position" p ON e."Position_ID" = p.id
		 WHERE e.id = $1`,
		userID,
	).Scan(
		&user.ID,
		&user.LastName,
		&user.FirstName,
		&user.MiddleName,
		&user.Phone,
		&user.DateOfBirth,
		&user.Email,
		&user.IsActive,
		&user.RoleID,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Получение данных пользователя (без пароля) по email
func (p *Provider) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User

	err := p.conn.QueryRow(
		`SELECT e.id, e."LastName", e."FirstName", e."MiddleName", e."Phone", e."BirthDate", e."Email", e."IsActive", p.id, p."Name"
		 FROM public."Employee" e
		 LEFT JOIN public."Position" p ON e."Position_ID" = p.id
		 WHERE lower(e."Email") = lower($1)`,
		email,
	).Scan(
		&user.ID,
		&user.LastName,
		&user.FirstName,
		&user.MiddleName,
		&user.Phone,
		&user.DateOfBirth,
		&user.Email,
		&user.IsActive,
		&user.RoleID,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Получение хеша пароля по почте
func (p *Provider) GetUserPasswordHashByEmail(email string) (string, error) {
	var passwordHash string
	err := p.conn.QueryRow(
		`SELECT "PasswordHash" FROM public."Employee" WHERE lower("Email") = lower($1) AND "IsActive" = TRUE`,
		email,
	).Scan(&passwordHash)
	return passwordHash, err
}

// Получение хеша пароля по id
func (p *Provider) GetUserPasswordHashByID(userID int) (string, error) {
	var passwordHash string
	err := p.conn.QueryRow(
		`SELECT "PasswordHash" FROM public."Employee" WHERE id = $1 AND "IsActive" = TRUE`,
		userID,
	).Scan(&passwordHash)
	return passwordHash, err
}

// Получение прав пользователя по id или почте
func (p *Provider) GetUserRoleByID(userID int) (entities.Role, error) {
	var role entities.Role
	err := p.conn.QueryRow(
		`SELECT p.id, TRIM(p."Name"), p."Description"
		 FROM public."Employee" e
		 JOIN public."Position" p ON p.id = e."Position_ID"
		 WHERE e.id = $1`,
		userID,
	).Scan(&role.ID, &role.Title, &role.Description)
	return role, err
}

func (p *Provider) GetUserRoleByEmail(email string) (entities.Role, error) {

	var role entities.Role
	err := p.conn.QueryRow(
		`SELECT p.id, TRIM(p."Name"), p."Description"
		 FROM public."Employee" e
		 JOIN public."Position" p ON p.id = e."Position_ID"
		 WHERE lower(e."Email") = lower($1)`,
		email,
	).Scan(&role.ID, &role.Title, &role.Description)
	return role, err
}

// Изменение хеша пароля пользователя
func (p *Provider) UpdateUserPasswordHash(userID int, newPasswordHash string) error {
	_, err := p.conn.Exec(
		`UPDATE public."Employee" SET "PasswordHash" = $2 WHERE id = $1`,
		userID, newPasswordHash,
	)
	return err
}

// Удаление пользователя (Изменяем IsActive на false)
func (p *Provider) ChangeUserActive(userID int, isActive bool) error {
	_, err := p.conn.Exec(
		`UPDATE public."Employee" SET "IsActive" = $2 WHERE id = $1`,
		userID, isActive,
	)
	return err
}

// Поиск пользователя по почте (фио, почта, права, телефон)
func (p *Provider) FindUsers(user entities.User, searchBy string) ([]entities.User, error) {
	var query string
	var searchValue string

	switch searchBy {
	case "email":
		query = `SELECT * FROM search_employee_flexible(p_email := $1)`
		searchValue = user.Email
	case "position":
		query = `SELECT * FROM search_employee_flexible(p_position_name := $1)`
		searchValue = user.Role
	case "last_name":
		query = `SELECT * FROM search_employee_flexible(p_last_name := $1)`
		searchValue = user.LastName
	case "first_name":
		query = `SELECT * FROM search_employee_flexible(p_first_name := $1)`
		searchValue = user.FirstName
	default:
		return []entities.User{}, nil
	}

	if searchValue == "" {
		return []entities.User{}, nil
	}

	rows, err := p.conn.Query(query, searchValue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entities.User

	for rows.Next() {
		var user entities.User
		// Порядок должен соответствовать RETURNS TABLE в функции:
		// id, last_name, first_name, middle_name, email, phone, birth_date, is_active, position_name, position_id
		err := rows.Scan(
			&user.ID,          // id
			&user.LastName,    // last_name
			&user.FirstName,   // first_name
			&user.MiddleName,  // middle_name
			&user.Email,       // email
			&user.Phone,       // phone
			&user.DateOfBirth, // birth_date
			&user.IsActive,    // is_active
			&user.Role,        // position_name - сохраняем в Role
			&user.RoleID,      // position_id
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

// Поиск пользователя по правам (Администратор, Руководитель, Редактор, Автор)
func (p *Provider) ListUsersByRole(roleID int) ([]entities.User, error) {
	rows, err := p.conn.Query(
		`SELECT e.id, e."LastName", e."FirstName", e."MiddleName", e."Phone", e."BirthDate", e."Email", e."IsActive", p.id, p."Name"
		 FROM public."Employee" e
		 JOIN public."Position" p ON p.id = e."Position_ID"
		 WHERE p.id = $1
		 ORDER BY e.id`,
		roleID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(
			&user.ID,
			&user.LastName,
			&user.FirstName,
			&user.MiddleName,
			&user.Phone,
			&user.DateOfBirth,
			&user.Email,
			&user.IsActive,
			&user.RoleID,
			&user.Role,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// Поиск пользователя по правам (Администратор, Руководитель, Редактор, Автор)
func (p *Provider) ListUsers() ([]entities.User, error) {
	rows, err := p.conn.Query(
		`SELECT e.id, e."LastName", e."FirstName", e."MiddleName", e."Phone", e."BirthDate", e."Email", e."IsActive", p.id, p."Name"
		 FROM public."Employee" e
		 JOIN public."Position" p ON p.id = e."Position_ID"
		 ORDER BY e.id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(
			&user.ID,
			&user.LastName,
			&user.FirstName,
			&user.MiddleName,
			&user.Phone,
			&user.DateOfBirth,
			&user.Email,
			&user.IsActive,
			&user.RoleID,
			&user.Role,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// Изменение прав пользователя (Администратор, Руководитель, Редактор, Автор)
func (p *Provider) UpdateUserRole(userID int, roleID int) error {
	_, err := p.conn.Exec(
		`UPDATE public."Employee" SET "Position_ID" = $2 WHERE id = $1`,
		userID, roleID,
	)
	return err
}

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
