package provider

import (
	"strings"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
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

	if strings.TrimSpace(searchValue) == "" {
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
