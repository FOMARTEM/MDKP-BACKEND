package api

import (
	"errors"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

// usecaseMock — минимальный мок интерфейса Usecase без внешних зависимостей.
// В тестах задавайте нужные Func поля; по умолчанию методы возвращают ошибку.
type usecaseMock struct {
	AuthorizeFunc          func(user entities.User) (*bool, error)
	SelectUserByEmailFunc  func(email string) (*entities.User, error)
	UpdateUserPasswordFunc func(userID int, newPassword string) error
	SelectUserByIDFunc     func(userID int) (*entities.User, error)
	GetLogsFunc            func(userID int, email, startDate, endDate string, limit, offset int) ([]entities.Log, error)
	GetRolesFunc           func() ([]entities.Role, error)
	GetStatusesFunc        func() ([]entities.Status, error)
	CreateUserFunc         func(user entities.User) (*entities.User, error)
	UserActiveChangeFunc   func(userEmail string) error
	GetUsersByRoleFunc     func(user entities.User) ([]entities.User, error)
	UserRoleUpdateFunc     func(userEmail string, roleID int) error
	GetUsersFunc           func() ([]entities.User, error)

	CreateTaskFunc       func(task entities.Tasks) (*entities.Tasks, error)
	TaskDeleteFunc       func(id int) error
	TaskGetByIdFunc      func(taskId int) (*entities.Tasks, error)
	TaskStatusUpdateFunc func(taskId int, status string) error
	TasksListFunc        func(userID int) ([]entities.Tasks, error)

	CreateMaterialFunc      func(material entities.Material) (int, error)
	GetMaterialFunc         func(materialId int) (*entities.Material, error)
	GetMaterialsByTaskIDFun func(taskID int) ([]entities.Material, error)

	VersionTaskFunc   func(version entities.Version) (*entities.Version, error)
	VersionsListFunc  func(taskId int) ([]entities.Version, error)
	VersionByIdFunc   func(versionId int) (*entities.Version, error)
	CreateRevisionFun func(revision entities.Revision) (*entities.Revision, error)
	EditStatusUpdateF func(editId int, status string) error
	GetRevisionsByVer func(versionID int) ([]entities.Revision, error)
	GetRevisionByIDF  func(revisionID int) (*entities.Revision, error)
}

var errMockNotImplemented = errors.New("mock: not implemented")

func boolPtr(v bool) *bool { return &v }

func (m *usecaseMock) Authorize(user entities.User) (*bool, error) {
	if m.AuthorizeFunc == nil {
		return boolPtr(false), errMockNotImplemented
	}
	return m.AuthorizeFunc(user)
}
func (m *usecaseMock) SelectUserByEmail(email string) (*entities.User, error) {
	if m.SelectUserByEmailFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.SelectUserByEmailFunc(email)
}
func (m *usecaseMock) UpdateUserPassword(userID int, newPassword string) error {
	if m.UpdateUserPasswordFunc == nil {
		return errMockNotImplemented
	}
	return m.UpdateUserPasswordFunc(userID, newPassword)
}
func (m *usecaseMock) SelectUserByID(userID int) (*entities.User, error) {
	if m.SelectUserByIDFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.SelectUserByIDFunc(userID)
}
func (m *usecaseMock) GetLogs(userID int, email, startDate, endDate string, limit, offset int) ([]entities.Log, error) {
	if m.GetLogsFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.GetLogsFunc(userID, email, startDate, endDate, limit, offset)
}
func (m *usecaseMock) GetRoles() ([]entities.Role, error) {
	if m.GetRolesFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.GetRolesFunc()
}
func (m *usecaseMock) GetStatuses() ([]entities.Status, error) {
	if m.GetStatusesFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.GetStatusesFunc()
}
func (m *usecaseMock) CreateUser(user entities.User) (*entities.User, error) {
	if m.CreateUserFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.CreateUserFunc(user)
}
func (m *usecaseMock) UserActiveChange(userEmail string) error {
	if m.UserActiveChangeFunc == nil {
		return errMockNotImplemented
	}
	return m.UserActiveChangeFunc(userEmail)
}
func (m *usecaseMock) GetUsersByRole(user entities.User) ([]entities.User, error) {
	if m.GetUsersByRoleFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.GetUsersByRoleFunc(user)
}
func (m *usecaseMock) UserRoleUpdate(userEmail string, roleID int) error {
	if m.UserRoleUpdateFunc == nil {
		return errMockNotImplemented
	}
	return m.UserRoleUpdateFunc(userEmail, roleID)
}
func (m *usecaseMock) GetUsers() ([]entities.User, error) {
	if m.GetUsersFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.GetUsersFunc()
}

func (m *usecaseMock) CreateTask(task entities.Tasks) (*entities.Tasks, error) {
	if m.CreateTaskFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.CreateTaskFunc(task)
}
func (m *usecaseMock) TaskDelete(id int) error {
	if m.TaskDeleteFunc == nil {
		return errMockNotImplemented
	}
	return m.TaskDeleteFunc(id)
}
func (m *usecaseMock) TaskGetById(taskId int) (*entities.Tasks, error) {
	if m.TaskGetByIdFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.TaskGetByIdFunc(taskId)
}
func (m *usecaseMock) TaskStatusUpdate(taskId int, status string) error {
	if m.TaskStatusUpdateFunc == nil {
		return errMockNotImplemented
	}
	return m.TaskStatusUpdateFunc(taskId, status)
}
func (m *usecaseMock) TasksList(userID int) ([]entities.Tasks, error) {
	if m.TasksListFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.TasksListFunc(userID)
}

func (m *usecaseMock) CreateMaterial(material entities.Material) (int, error) {
	if m.CreateMaterialFunc == nil {
		return 0, errMockNotImplemented
	}
	return m.CreateMaterialFunc(material)
}
func (m *usecaseMock) GetMaterial(materialId int) (*entities.Material, error) {
	if m.GetMaterialFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.GetMaterialFunc(materialId)
}
func (m *usecaseMock) GetMaterialsByTaskID(taskID int) ([]entities.Material, error) {
	if m.GetMaterialsByTaskIDFun == nil {
		return nil, errMockNotImplemented
	}
	return m.GetMaterialsByTaskIDFun(taskID)
}

func (m *usecaseMock) VersionTask(version entities.Version) (*entities.Version, error) {
	if m.VersionTaskFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.VersionTaskFunc(version)
}
func (m *usecaseMock) VersionsList(taskId int) ([]entities.Version, error) {
	if m.VersionsListFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.VersionsListFunc(taskId)
}
func (m *usecaseMock) VersionById(versionId int) (*entities.Version, error) {
	if m.VersionByIdFunc == nil {
		return nil, errMockNotImplemented
	}
	return m.VersionByIdFunc(versionId)
}
func (m *usecaseMock) CreateRevision(revision entities.Revision) (*entities.Revision, error) {
	if m.CreateRevisionFun == nil {
		return nil, errMockNotImplemented
	}
	return m.CreateRevisionFun(revision)
}
func (m *usecaseMock) EditStatusUpdate(editId int, status string) error {
	if m.EditStatusUpdateF == nil {
		return errMockNotImplemented
	}
	return m.EditStatusUpdateF(editId, status)
}
func (m *usecaseMock) GetRevisionsByVersionID(versionID int) ([]entities.Revision, error) {
	if m.GetRevisionsByVer == nil {
		return nil, errMockNotImplemented
	}
	return m.GetRevisionsByVer(versionID)
}
func (m *usecaseMock) GetRevisionByID(revisionID int) (*entities.Revision, error) {
	if m.GetRevisionByIDF == nil {
		return nil, errMockNotImplemented
	}
	return m.GetRevisionByIDF(revisionID)
}

