package usecase

import (
	"errors"
	"testing"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

func TestGetUsersByRole_EmailLastFirstBranchesAndProviderErr(t *testing.T) {
	t.Parallel()

	dbErr := errors.New("find err")
	p := &providerMock{
		FindUsersFunc: func(user entities.User, searchBy string) ([]entities.User, error) {
			if user.Email != "" && searchBy != "email" {
				t.Fatalf("expected searchBy=email, got %q", searchBy)
			}
			if user.LastName != "" && searchBy != "last_name" {
				t.Fatalf("expected searchBy=last_name, got %q", searchBy)
			}
			if user.FirstName != "" && searchBy != "first_name" {
				t.Fatalf("expected searchBy=first_name, got %q", searchBy)
			}
			if user.FirstName == "ERR" {
				return nil, dbErr
			}
			return []entities.User{}, nil
		},
	}
	u := NewUsecase(p)

	if _, err := u.GetUsersByRole(entities.User{Email: "a@b.c"}); err != nil {
		t.Fatalf("email branch err: %v", err)
	}
	if _, err := u.GetUsersByRole(entities.User{LastName: "L"}); err != nil {
		t.Fatalf("last_name branch err: %v", err)
	}
	if _, err := u.GetUsersByRole(entities.User{FirstName: "F"}); err != nil {
		t.Fatalf("first_name branch err: %v", err)
	}
	if _, err := u.GetUsersByRole(entities.User{FirstName: "ERR"}); err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected dbErr, got %v", err)
	}
}

func TestUserRoleUpdate_OK_AndErr(t *testing.T) {
	t.Parallel()

	dbErr := errors.New("update err")
	p := &providerMock{
		GetUserByEmailFunc: func(email string) (*entities.User, error) { return &entities.User{ID: 9}, nil },
		UpdateUserRoleFunc: func(userID int, roleID int) error {
			if userID != 9 || roleID != 2 {
				t.Fatalf("unexpected args: %d %d", userID, roleID)
			}
			return nil
		},
	}
	u := NewUsecase(p)
	if err := u.UserRoleUpdate("a@b.c", 2); err != nil {
		t.Fatalf("UserRoleUpdate: %v", err)
	}

	p2 := &providerMock{
		GetUserByEmailFunc: func(email string) (*entities.User, error) { return &entities.User{ID: 9}, nil },
		UpdateUserRoleFunc: func(userID int, roleID int) error { return dbErr },
	}
	u2 := NewUsecase(p2)
	if err := u2.UserRoleUpdate("a@b.c", 2); err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected dbErr, got %v", err)
	}
}

func TestGetUsers_OK_AndErr(t *testing.T) {
	t.Parallel()

	dbErr := errors.New("list err")
	u := NewUsecase(&providerMock{
		ListUsersFunc: func() ([]entities.User, error) { return []entities.User{{ID: 1}}, nil },
	})
	users, err := u.GetUsers()
	if err != nil || len(users) != 1 || users[0].ID != 1 {
		t.Fatalf("unexpected: users=%#v err=%v", users, err)
	}

	u2 := NewUsecase(&providerMock{
		ListUsersFunc: func() ([]entities.User, error) { return nil, dbErr },
	})
	_, err = u2.GetUsers()
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected dbErr, got %v", err)
	}
}

func TestTasksMaterialsVersionsRevisions(t *testing.T) {
	t.Parallel()

	p := &providerMock{
		CreateTaskFunc: func(task entities.Tasks) (*entities.Tasks, error) {
			task.ID = 1
			return &task, nil
		},
		DeleteTaskFunc:       func(id int) error { return nil },
		GetTaskByIDFunc:      func(taskID int) (*entities.Tasks, error) { return &entities.Tasks{ID: taskID}, nil },
		UpdateTaskStatusFunc: func(taskID int, status string) error { return nil },
		GetTasksByUserIDFunc: func(userID int) ([]entities.Tasks, error) { return []entities.Tasks{{ID: 2}}, nil },

		CreateMaterialFunc:       func(material entities.Material) (*entities.Material, error) { material.ID = 7; return &material, nil },
		GetMaterialByIDFunc:      func(materialID int) (*entities.Material, error) { return &entities.Material{ID: materialID}, nil },
		GetMaterialsByTaskIDFunc: func(taskID int) ([]entities.Material, error) { return []entities.Material{{ID: 3}}, nil },

		CreateVersionFunc:       func(version entities.Version) (*entities.Version, error) { version.ID = 11; return &version, nil },
		GetVersionsByTaskIDFunc: func(taskID int) ([]entities.Version, error) { return []entities.Version{{ID: 4}}, nil },
		GetVersionByIDFunc:      func(versionID int) (*entities.Version, error) { return &entities.Version{ID: versionID}, nil },

		CreateRevisionFunc:       func(revision entities.Revision) (*entities.Revision, error) { revision.ID = 13; return &revision, nil },
		UpdateRevisionStatusFunc: func(revisionID int, status string) error { return nil },
		GetRevisionsByVersionIDF: func(versionID int) ([]entities.Revision, error) { return []entities.Revision{{ID: 5}}, nil },
		GetRevisionByIDFunc:      func(revisionID int) (*entities.Revision, error) { return &entities.Revision{ID: revisionID}, nil },

		GetRolesFunc:    func() ([]entities.Role, error) { return []entities.Role{{ID: 1}}, nil },
		GetStatusesFunc: func() ([]entities.Status, error) { return []entities.Status{{ID: 1}}, nil },
	}
	u := NewUsecase(p)

	if _, err := u.CreateTask(entities.Tasks{Title: "t", Description: "d", Priority: 1}); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if err := u.TaskDelete(1); err != nil {
		t.Fatalf("TaskDelete: %v", err)
	}
	if _, err := u.TaskGetById(1); err != nil {
		t.Fatalf("TaskGetById: %v", err)
	}
	if err := u.TaskStatusUpdate(1, "x"); err != nil {
		t.Fatalf("TaskStatusUpdate: %v", err)
	}
	if _, err := u.TasksList(1); err != nil {
		t.Fatalf("TasksList: %v", err)
	}

	if id, err := u.CreateMaterial(entities.Material{}); err != nil || id != 7 {
		t.Fatalf("CreateMaterial: id=%d err=%v", id, err)
	}
	if _, err := u.GetMaterial(7); err != nil {
		t.Fatalf("GetMaterial: %v", err)
	}
	if _, err := u.GetMaterialsByTaskID(1); err != nil {
		t.Fatalf("GetMaterialsByTaskID: %v", err)
	}

	if _, err := u.VersionTask(entities.Version{}); err != nil {
		t.Fatalf("VersionTask: %v", err)
	}
	if _, err := u.VersionsList(1); err != nil {
		t.Fatalf("VersionsList: %v", err)
	}
	if _, err := u.VersionById(1); err != nil {
		t.Fatalf("VersionById: %v", err)
	}

	if _, err := u.CreateRevision(entities.Revision{}); err != nil {
		t.Fatalf("CreateRevision: %v", err)
	}
	if err := u.EditStatusUpdate(1, "x"); err != nil {
		t.Fatalf("EditStatusUpdate: %v", err)
	}
	if _, err := u.GetRevisionsByVersionID(1); err != nil {
		t.Fatalf("GetRevisionsByVersionID: %v", err)
	}
	if _, err := u.GetRevisionByID(1); err != nil {
		t.Fatalf("GetRevisionByID: %v", err)
	}

	if _, err := u.GetRoles(); err != nil {
		t.Fatalf("GetRoles: %v", err)
	}
	if _, err := u.GetStatuses(); err != nil {
		t.Fatalf("GetStatuses: %v", err)
	}
}

