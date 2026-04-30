package usecase

import "github.com/FOMARTEM/MDKP-BACKEND/internal/entities"

type Provider interface {
	// Создание пользователя
	CreateUser(user entities.User) (*entities.User, error)

	// Получение данных пользователя (без пароля) по id
	GetUserByID(userID int) (*entities.User, error)

	// Получение данных пользователя (без пароля) по email
	GetUserByEmail(email string) (*entities.User, error)

	// Получение хеша пароля по почте
	GetUserPasswordHashByEmail(email string) (string, error)

	// Получение хеша пароля по id
	GetUserPasswordHashByID(userID int) (string, error)

	// Получение прав пользователя по id или почте
	GetUserRoleByID(userID int) (entities.Role, error)
	GetUserRoleByEmail(email string) (entities.Role, error)

	// Изменение хеша пароля пользователя
	UpdateUserPasswordHash(userID int, newPasswordHash string) error

	// Удаление пользователя
	ChangeUserActive(userID int, isActive bool) error

	// Поиск пользователя по почте (фио, почта, права, телефон)
	FindUsers(user entities.User, searchBy string) ([]entities.User, error)

	// Поиск пользователя по правам (Администратор, Руководитель, Редактор, Автор)
	ListUsersByRole(roleID int) ([]entities.User, error)

	// Список всех пользователей
	ListUsers() ([]entities.User, error)

	// Изменение прав пользователя (Администратор, Руководитель, Редактор, Автор)
	UpdateUserRole(userID int, roleID int) error

	// Создание лога (При любом действии пользователя, при входе в систему, при выходе из системы, при изменении данных аккаунта, при создании задачи, при удалении задачи, при изменении статуса задачи, при назначении задачи, при создании правки, при удалении правки, при изменении статуса правки)
	CreateLog(userID int, action string) error

	// Получение логов по id пользователя
	GetLogsByUserID(userID int) ([]entities.Log, error)

	// Получение логов по почте пользователя
	GetLogsByUserEmail(email string) ([]entities.Log, error)

	// Получение логов по определённым датам
	GetLogsByDateRange(startDate, endDate string) ([]entities.Log, error)

	// Создание задачи
	CreateTask(task entities.Tasks) (*entities.Tasks, error)

	// Удаление задачи
	DeleteTask(id int) error

	// Получение списка задач пользователя (по id)
	GetTasksByUserID(userID int) ([]entities.Tasks, error)

	// Получение задачи по id
	GetTaskByID(taskID int) (*entities.Tasks, error)

	// Изменение статуса задачи (Открыта, В работе, На проверке, Закрыта)
	UpdateTaskStatus(taskID int, status string) error

	// Изменение автора задачи
	UpdateTaskAuthor(taskID int, authorID int, authorEmail string) error

	// Изменение редактора задачи
	UpdateTaskEditor(taskID int, editorID int, editorEmail string) error

	// Установка даты готовности задачи
	UpdateTaskReadyDate(taskID int, readyDate string) error

	// Получение задач по приоритетам
	GetTasksByPriority(priority int) ([]entities.Tasks, error)

	// Получение задач с определёнными статусами
	GetTasksByStatuses(statuses []string) ([]entities.Tasks, error)

	// Создание версии задачи
	CreateVersion(version entities.Version) (*entities.Version, error)

	// Получение версий задачи по id задачи
	GetVersionsByTaskID(taskID int) ([]entities.Version, error)

	// Получение версии задачи по id версии
	GetVersionByID(versionID int) (*entities.Version, error)

	// Изменение статуса версии (Открыта, В работе, На проверке, Закрыта)
	UpdateVersionStatus(versionID int, status string) error

	// Создание материала
	CreateMaterial(material entities.Material) (*entities.Material, error)

	// Получение материалов по id задачи
	GetMaterialsByTaskID(taskID int) ([]entities.Material, error)

	// Получение материала по id материала
	GetMaterialByID(materialID int) (*entities.Material, error)

	// Присвоение материала к задаче
	AssignMaterialToTask(materialID int, taskID int) error

	// Присвоение материала к версии
	AssignMaterialToVersion(materialID int, versionID int) error

	// Создание правки
	CreateRevision(revision entities.Revision) (*entities.Revision, error)

	// Получение правок по id задачи
	GetRevisionsByTaskID(taskID int) ([]entities.Revision, error)

	// Получение правки по id правки
	GetRevisionByID(revisionID int) (*entities.Revision, error)

	// Изменение статуса правки (Открыта, В работе, На проверке, Закрыта)
	UpdateRevisionStatus(revisionID int, status string) error

	// Получение ролей пользователей
	GetRoles() ([]entities.Role, error)

	// Получение Id роли по её названию
	GetRoleId(role string) (int, error)

	GetStatuses() ([]entities.Status, error)
}
