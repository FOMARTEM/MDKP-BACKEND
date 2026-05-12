package usecase

import (
	"errors"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

// providerMock реализует интерфейс Provider через набор Func полей.
// В тестах задавайте только то, что нужно; по умолчанию методы возвращают ошибку.
type providerMock struct {
	CreateUserFunc func(user entities.User) (*entities.User, error)

	GetUserByIDFunc               func(userID int) (*entities.User, error)
	GetUserByEmailFunc            func(email string) (*entities.User, error)
	GetUserPasswordHashByEmailFun func(email string) (string, error)
	GetUserPasswordHashByIDFunc   func(userID int) (string, error)

	GetUserRoleByIDFunc    func(userID int) (entities.Role, error)
	GetUserRoleByEmailFunc func(email string) (entities.Role, error)

	UpdateUserPasswordHashFunc func(userID int, newPasswordHash string) error
	ChangeUserActiveFunc       func(userID int, isActive bool) error
	FindUsersFunc              func(user entities.User, searchBy string) ([]entities.User, error)
	ListUsersByRoleFunc        func(roleID int) ([]entities.User, error)
	ListUsersFunc              func() ([]entities.User, error)
	UpdateUserRoleFunc         func(userID int, roleID int) error

	CreateLogFunc          func(userID int, action string) error
	GetLogsByUserIDFunc    func(userID int) ([]entities.Log, error)
	GetLogsByUserEmailFunc func(email string) ([]entities.Log, error)
	GetLogsByDateRangeFunc func(startDate, endDate string) ([]entities.Log, error)
	GetLogsAllFunc         func(limit, offset int) ([]entities.Log, error)

	CreateTaskFunc          func(task entities.Tasks) (*entities.Tasks, error)
	DeleteTaskFunc          func(id int) error
	GetTasksByUserIDFunc    func(userID int) ([]entities.Tasks, error)
	GetTaskByIDFunc         func(taskID int) (*entities.Tasks, error)
	UpdateTaskStatusFunc    func(taskID int, status string) error
	UpdateTaskAuthorFunc    func(taskID int, authorID int, authorEmail string) error
	UpdateTaskEditorFunc    func(taskID int, editorID int, editorEmail string) error
	UpdateTaskReadyDateFunc func(taskID int, readyDate string) error
	GetTasksByPriorityFunc  func(priority int) ([]entities.Tasks, error)
	GetTasksByStatusesFunc  func(statuses []string) ([]entities.Tasks, error)

	CreateVersionFunc       func(version entities.Version) (*entities.Version, error)
	GetVersionsByTaskIDFunc func(taskID int) ([]entities.Version, error)
	GetVersionByIDFunc      func(versionID int) (*entities.Version, error)
	UpdateVersionStatusFunc func(versionID int, status string) error

	CreateMaterialFunc         func(material entities.Material) (*entities.Material, error)
	GetMaterialsByTaskIDFunc   func(taskID int) ([]entities.Material, error)
	GetMaterialByIDFunc        func(materialID int) (*entities.Material, error)
	AssignMaterialToTaskFunc   func(materialID int, taskID int) error
	AssignMaterialToVersionFun func(materialID int, versionID int) error

	CreateRevisionFunc       func(revision entities.Revision) (*entities.Revision, error)
	GetRevisionsByVersionIDF func(versionID int) ([]entities.Revision, error)
	GetRevisionByIDFunc      func(revisionID int) (*entities.Revision, error)
	UpdateRevisionStatusFunc func(revisionID int, status string) error

	GetRolesFunc    func() ([]entities.Role, error)
	GetRoleIdFunc   func(role string) (int, error)
	GetStatusesFunc func() ([]entities.Status, error)
}

var errProviderMockNotImplemented = errors.New("provider mock: not implemented")

func (m *providerMock) CreateUser(user entities.User) (*entities.User, error) {
	if m.CreateUserFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.CreateUserFunc(user)
}

func (m *providerMock) GetUserByID(userID int) (*entities.User, error) {
	if m.GetUserByIDFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetUserByIDFunc(userID)
}
func (m *providerMock) GetUserByEmail(email string) (*entities.User, error) {
	if m.GetUserByEmailFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetUserByEmailFunc(email)
}
func (m *providerMock) GetUserPasswordHashByEmail(email string) (string, error) {
	if m.GetUserPasswordHashByEmailFun == nil {
		return "", errProviderMockNotImplemented
	}
	return m.GetUserPasswordHashByEmailFun(email)
}
func (m *providerMock) GetUserPasswordHashByID(userID int) (string, error) {
	if m.GetUserPasswordHashByIDFunc == nil {
		return "", errProviderMockNotImplemented
	}
	return m.GetUserPasswordHashByIDFunc(userID)
}
func (m *providerMock) GetUserRoleByID(userID int) (entities.Role, error) {
	if m.GetUserRoleByIDFunc == nil {
		return entities.Role{}, errProviderMockNotImplemented
	}
	return m.GetUserRoleByIDFunc(userID)
}
func (m *providerMock) GetUserRoleByEmail(email string) (entities.Role, error) {
	if m.GetUserRoleByEmailFunc == nil {
		return entities.Role{}, errProviderMockNotImplemented
	}
	return m.GetUserRoleByEmailFunc(email)
}
func (m *providerMock) UpdateUserPasswordHash(userID int, newPasswordHash string) error {
	if m.UpdateUserPasswordHashFunc == nil {
		return errProviderMockNotImplemented
	}
	return m.UpdateUserPasswordHashFunc(userID, newPasswordHash)
}
func (m *providerMock) ChangeUserActive(userID int, isActive bool) error {
	if m.ChangeUserActiveFunc == nil {
		return errProviderMockNotImplemented
	}
	return m.ChangeUserActiveFunc(userID, isActive)
}
func (m *providerMock) FindUsers(user entities.User, searchBy string) ([]entities.User, error) {
	if m.FindUsersFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.FindUsersFunc(user, searchBy)
}
func (m *providerMock) ListUsersByRole(roleID int) ([]entities.User, error) {
	if m.ListUsersByRoleFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.ListUsersByRoleFunc(roleID)
}
func (m *providerMock) ListUsers() ([]entities.User, error) {
	if m.ListUsersFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.ListUsersFunc()
}
func (m *providerMock) UpdateUserRole(userID int, roleID int) error {
	if m.UpdateUserRoleFunc == nil {
		return errProviderMockNotImplemented
	}
	return m.UpdateUserRoleFunc(userID, roleID)
}

func (m *providerMock) CreateLog(userID int, action string) error {
	if m.CreateLogFunc == nil {
		return errProviderMockNotImplemented
	}
	return m.CreateLogFunc(userID, action)
}
func (m *providerMock) GetLogsByUserID(userID int) ([]entities.Log, error) {
	if m.GetLogsByUserIDFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetLogsByUserIDFunc(userID)
}
func (m *providerMock) GetLogsByUserEmail(email string) ([]entities.Log, error) {
	if m.GetLogsByUserEmailFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetLogsByUserEmailFunc(email)
}
func (m *providerMock) GetLogsByDateRange(startDate, endDate string) ([]entities.Log, error) {
	if m.GetLogsByDateRangeFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetLogsByDateRangeFunc(startDate, endDate)
}
func (m *providerMock) GetLogsAll(limit, offset int) ([]entities.Log, error) {
	if m.GetLogsAllFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetLogsAllFunc(limit, offset)
}

func (m *providerMock) CreateTask(task entities.Tasks) (*entities.Tasks, error) {
	if m.CreateTaskFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.CreateTaskFunc(task)
}
func (m *providerMock) DeleteTask(id int) error {
	if m.DeleteTaskFunc == nil {
		return errProviderMockNotImplemented
	}
	return m.DeleteTaskFunc(id)
}
func (m *providerMock) GetTasksByUserID(userID int) ([]entities.Tasks, error) {
	if m.GetTasksByUserIDFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetTasksByUserIDFunc(userID)
}
func (m *providerMock) GetTaskByID(taskID int) (*entities.Tasks, error) {
	if m.GetTaskByIDFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetTaskByIDFunc(taskID)
}
func (m *providerMock) UpdateTaskStatus(taskID int, status string) error {
	if m.UpdateTaskStatusFunc == nil {
		return errProviderMockNotImplemented
	}
	return m.UpdateTaskStatusFunc(taskID, status)
}
func (m *providerMock) UpdateTaskAuthor(taskID int, authorID int, authorEmail string) error {
	if m.UpdateTaskAuthorFunc == nil {
		return errProviderMockNotImplemented
	}
	return m.UpdateTaskAuthorFunc(taskID, authorID, authorEmail)
}
func (m *providerMock) UpdateTaskEditor(taskID int, editorID int, editorEmail string) error {
	if m.UpdateTaskEditorFunc == nil {
		return errProviderMockNotImplemented
	}
	return m.UpdateTaskEditorFunc(taskID, editorID, editorEmail)
}
func (m *providerMock) UpdateTaskReadyDate(taskID int, readyDate string) error {
	if m.UpdateTaskReadyDateFunc == nil {
		return errProviderMockNotImplemented
	}
	return m.UpdateTaskReadyDateFunc(taskID, readyDate)
}
func (m *providerMock) GetTasksByPriority(priority int) ([]entities.Tasks, error) {
	if m.GetTasksByPriorityFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetTasksByPriorityFunc(priority)
}
func (m *providerMock) GetTasksByStatuses(statuses []string) ([]entities.Tasks, error) {
	if m.GetTasksByStatusesFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetTasksByStatusesFunc(statuses)
}

func (m *providerMock) CreateVersion(version entities.Version) (*entities.Version, error) {
	if m.CreateVersionFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.CreateVersionFunc(version)
}
func (m *providerMock) GetVersionsByTaskID(taskID int) ([]entities.Version, error) {
	if m.GetVersionsByTaskIDFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetVersionsByTaskIDFunc(taskID)
}
func (m *providerMock) GetVersionByID(versionID int) (*entities.Version, error) {
	if m.GetVersionByIDFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetVersionByIDFunc(versionID)
}
func (m *providerMock) UpdateVersionStatus(versionID int, status string) error {
	if m.UpdateVersionStatusFunc == nil {
		return errProviderMockNotImplemented
	}
	return m.UpdateVersionStatusFunc(versionID, status)
}

func (m *providerMock) CreateMaterial(material entities.Material) (*entities.Material, error) {
	if m.CreateMaterialFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.CreateMaterialFunc(material)
}
func (m *providerMock) GetMaterialsByTaskID(taskID int) ([]entities.Material, error) {
	if m.GetMaterialsByTaskIDFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetMaterialsByTaskIDFunc(taskID)
}
func (m *providerMock) GetMaterialByID(materialID int) (*entities.Material, error) {
	if m.GetMaterialByIDFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetMaterialByIDFunc(materialID)
}
func (m *providerMock) AssignMaterialToTask(materialID int, taskID int) error {
	if m.AssignMaterialToTaskFunc == nil {
		return errProviderMockNotImplemented
	}
	return m.AssignMaterialToTaskFunc(materialID, taskID)
}
func (m *providerMock) AssignMaterialToVersion(materialID int, versionID int) error {
	if m.AssignMaterialToVersionFun == nil {
		return errProviderMockNotImplemented
	}
	return m.AssignMaterialToVersionFun(materialID, versionID)
}

func (m *providerMock) CreateRevision(revision entities.Revision) (*entities.Revision, error) {
	if m.CreateRevisionFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.CreateRevisionFunc(revision)
}
func (m *providerMock) GetRevisionsByVersionID(versionID int) ([]entities.Revision, error) {
	if m.GetRevisionsByVersionIDF == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetRevisionsByVersionIDF(versionID)
}
func (m *providerMock) GetRevisionByID(revisionID int) (*entities.Revision, error) {
	if m.GetRevisionByIDFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetRevisionByIDFunc(revisionID)
}
func (m *providerMock) UpdateRevisionStatus(revisionID int, status string) error {
	if m.UpdateRevisionStatusFunc == nil {
		return errProviderMockNotImplemented
	}
	return m.UpdateRevisionStatusFunc(revisionID, status)
}

func (m *providerMock) GetRoles() ([]entities.Role, error) {
	if m.GetRolesFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetRolesFunc()
}
func (m *providerMock) GetRoleId(role string) (int, error) {
	if m.GetRoleIdFunc == nil {
		return 0, errProviderMockNotImplemented
	}
	return m.GetRoleIdFunc(role)
}
func (m *providerMock) GetStatuses() ([]entities.Status, error) {
	if m.GetStatusesFunc == nil {
		return nil, errProviderMockNotImplemented
	}
	return m.GetStatusesFunc()
}

