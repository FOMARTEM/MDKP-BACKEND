package usecase

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

func TestNewUsecase(t *testing.T) {
	t.Parallel()

	u := NewUsecase(&providerMock{})
	if u == nil {
		t.Fatalf("expected non-nil")
	}
}

func TestSelectUserByEmail_OK(t *testing.T) {
	t.Parallel()

	p := &providerMock{
		GetUserByEmailFunc: func(email string) (*entities.User, error) {
			return &entities.User{ID: 1, Email: email}, nil
		},
	}
	u := NewUsecase(p)

	got, err := u.SelectUserByEmail("a@b.c")
	if err != nil {
		t.Fatalf("SelectUserByEmail: %v", err)
	}
	if got == nil || got.Email != "a@b.c" {
		t.Fatalf("unexpected user: %#v", got)
	}
}

func TestSelectUserByEmail_NoRows(t *testing.T) {
	t.Parallel()

	p := &providerMock{
		GetUserByEmailFunc: func(email string) (*entities.User, error) {
			return nil, sql.ErrNoRows
		},
	}
	u := NewUsecase(p)

	_, err := u.SelectUserByEmail("a@b.c")
	if err == nil || !errors.Is(err, entities.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestSelectUserByEmail_OtherErr(t *testing.T) {
	t.Parallel()

	dbErr := errors.New("db")
	p := &providerMock{
		GetUserByEmailFunc: func(email string) (*entities.User, error) { return nil, dbErr },
	}
	u := NewUsecase(p)
	_, err := u.SelectUserByEmail("a@b.c")
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestSelectUserByID_NoRows(t *testing.T) {
	t.Parallel()

	p := &providerMock{
		GetUserByIDFunc: func(userID int) (*entities.User, error) { return nil, sql.ErrNoRows },
	}
	u := NewUsecase(p)
	_, err := u.SelectUserByID(1)
	if err == nil || !errors.Is(err, entities.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestSelectUserByID_OK(t *testing.T) {
	t.Parallel()

	p := &providerMock{
		GetUserByIDFunc: func(userID int) (*entities.User, error) { return &entities.User{ID: userID}, nil },
	}
	u := NewUsecase(p)
	got, err := u.SelectUserByID(7)
	if err != nil {
		t.Fatalf("SelectUserByID: %v", err)
	}
	if got == nil || got.ID != 7 {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestUpdateUserPassword_OK(t *testing.T) {
	t.Parallel()

	var gotHash string
	var gotLog string
	p := &providerMock{
		UpdateUserPasswordHashFunc: func(userID int, newPasswordHash string) error {
			if userID != 5 {
				t.Fatalf("unexpected userID: %d", userID)
			}
			gotHash = newPasswordHash
			return nil
		},
		CreateLogFunc: func(userID int, action string) error {
			gotLog = action
			return nil
		},
	}
	u := NewUsecase(p)

	if err := u.UpdateUserPassword(5, "Test12345"); err != nil {
		t.Fatalf("UpdateUserPassword: %v", err)
	}
	if gotHash == "" {
		t.Fatalf("expected password hash")
	}
	if gotLog != "Password updated" {
		t.Fatalf("unexpected log: %q", gotLog)
	}
}

func TestUpdateUserPassword_UpdateHashErr(t *testing.T) {
	t.Parallel()

	dbErr := errors.New("update err")
	p := &providerMock{
		UpdateUserPasswordHashFunc: func(userID int, newPasswordHash string) error { return dbErr },
	}
	u := NewUsecase(p)

	err := u.UpdateUserPassword(1, "Test12345")
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected dbErr, got %v", err)
	}
}

func TestUpdateUserPassword_LogErr(t *testing.T) {
	t.Parallel()

	dbErr := errors.New("log err")
	p := &providerMock{
		UpdateUserPasswordHashFunc: func(userID int, newPasswordHash string) error { return nil },
		CreateLogFunc:             func(userID int, action string) error { return dbErr },
	}
	u := NewUsecase(p)

	err := u.UpdateUserPassword(1, "Test12345")
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected dbErr, got %v", err)
	}
}

func TestGetLogs_AllBranchesAndPagination(t *testing.T) {
	t.Parallel()

	p := &providerMock{
		GetLogsByUserIDFunc: func(userID int) ([]entities.Log, error) {
			return []entities.Log{{ID: 1}, {ID: 2}, {ID: 3}}, nil
		},
		GetLogsByUserEmailFunc: func(email string) ([]entities.Log, error) {
			return []entities.Log{{ID: 4}, {ID: 5}}, nil
		},
		GetLogsByDateRangeFunc: func(startDate, endDate string) ([]entities.Log, error) {
			return []entities.Log{{ID: 6}}, nil
		},
		GetLogsAllFunc: func(limit, offset int) ([]entities.Log, error) {
			return []entities.Log{{ID: 7}}, nil
		},
	}
	u := NewUsecase(p)

	// userID branch with pagination
	logs, err := u.GetLogs(1, "", "", "", 2, 1)
	if err != nil || len(logs) != 2 || logs[0].ID != 2 || logs[1].ID != 3 {
		t.Fatalf("unexpected userID logs: %#v err=%v", logs, err)
	}

	// email branch, offset > len => []
	logs, err = u.GetLogs(0, "x@y.z", "", "", 10, 5)
	if err != nil || len(logs) != 0 {
		t.Fatalf("unexpected email logs: %#v err=%v", logs, err)
	}

	// date branch, end clamp
	logs, err = u.GetLogs(0, "", "2026-01-01", "2026-01-02", 10, 0)
	if err != nil || len(logs) != 1 || logs[0].ID != 6 {
		t.Fatalf("unexpected date logs: %#v err=%v", logs, err)
	}

	// no filters -> GetLogsAll direct
	logs, err = u.GetLogs(0, "", "", "", 10, 0)
	if err != nil || len(logs) != 1 || logs[0].ID != 7 {
		t.Fatalf("unexpected all logs: %#v err=%v", logs, err)
	}
}

func TestCreateUser_OK(t *testing.T) {
	t.Parallel()

	p := &providerMock{
		GetRoleIdFunc: func(role string) (int, error) {
			if role != "Админ" {
				t.Fatalf("unexpected role: %q", role)
			}
			return 2, nil
		},
		CreateUserFunc: func(user entities.User) (*entities.User, error) {
			if !user.IsActive || user.RoleID != 2 || user.PasswordHash == "" {
				t.Fatalf("unexpected user passed: %#v", user)
			}
			u := user
			u.ID = 10
			return &u, nil
		},
	}
	u := NewUsecase(p)

	created, err := u.CreateUser(entities.User{Email: "a@b.c", Password: "Test12345", Role: "Админ"})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if created == nil || created.ID != 10 {
		t.Fatalf("unexpected: %#v", created)
	}
}

func TestCreateUser_RoleErr(t *testing.T) {
	t.Parallel()

	dbErr := errors.New("role err")
	p := &providerMock{
		GetRoleIdFunc: func(role string) (int, error) { return 0, dbErr },
	}
	u := NewUsecase(p)
	_, err := u.CreateUser(entities.User{Password: "Test12345", Role: "Админ"})
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected dbErr, got %v", err)
	}
}

func TestUserActiveChange_OK(t *testing.T) {
	t.Parallel()

	var gotActive bool
	p := &providerMock{
		GetUserByEmailFunc: func(email string) (*entities.User, error) {
			return &entities.User{ID: 5, IsActive: true}, nil
		},
		ChangeUserActiveFunc: func(userID int, isActive bool) error {
			gotActive = isActive
			return nil
		},
	}
	u := NewUsecase(p)
	if err := u.UserActiveChange("a@b.c"); err != nil {
		t.Fatalf("UserActiveChange: %v", err)
	}
	if gotActive != false {
		t.Fatalf("expected toggle to false")
	}
}

func TestGetUsersByRole_SearchByAndNoCriteria(t *testing.T) {
	t.Parallel()

	var gotSearchBy string
	p := &providerMock{
		FindUsersFunc: func(user entities.User, searchBy string) ([]entities.User, error) {
			gotSearchBy = searchBy
			return []entities.User{{ID: 1}}, nil
		},
	}
	u := NewUsecase(p)

	_, err := u.GetUsersByRole(entities.User{Role: "Админ"})
	if err != nil {
		t.Fatalf("GetUsersByRole: %v", err)
	}
	if gotSearchBy != "position" {
		t.Fatalf("unexpected searchBy: %q", gotSearchBy)
	}

	_, err = u.GetUsersByRole(entities.User{})
	if err == nil {
		t.Fatalf("expected error for no criteria")
	}
}

