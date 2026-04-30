package api

import "github.com/FOMARTEM/MDKP-BACKEND/internal/entities"

type Usecase interface {
	// Пользователь
	// Функция для авторизации пользователя
	Authorize(user entities.User) (*bool, error)
	// Функция для получения данных пользователя (без пароля) по email
	SelectUserByEmail(email string) (*entities.User, error)
	// Функция для обновления пароля пользователя
	UpdateUserPassword(userID int, newPassword string) error
	// Функция для получения данных пользователя (без пароля) по id
	SelectUserByID(userID int) (*entities.User, error)
	// Получение логов
	GetLogs(userID int, email string, startDate, endDate string, limit, offset int) ([]entities.Log, error)
	// Получение ролей пользователей
	GetRoles() ([]entities.Role, error)
	// Получение возможных статусов
	GetStatuses() ([]entities.Status, error)
	// Создание пользователя
	CreateUser(user entities.User) (*entities.User, error)
	// Удаление пользователя
	UserActiveChange(userEmail string) error
	// Получаем всех пользователей с определёнными правами
	GetUsersByRole(user entities.User) ([]entities.User, error)
	// Обновление роли пользователя
	UserRoleUpdate(userEmail string, roleID int) error
	// Список пользователей
	GetUsers() ([]entities.User, error)

	// Работа с задачами
	// Создание задачи
	CreateTask(task entities.Tasks) (*entities.Tasks, error)
	// Удаление задачи
	TaskDelete(id int) error
	// Получение задачи по её Id
	TaskGetById(taskId int) (*entities.Tasks, error)
	// Изменение статуса задачи
	TaskStatusUpdate(taskId int, status string) error
}
