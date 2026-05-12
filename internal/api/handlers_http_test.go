package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

var chdirMu sync.Mutex

func mustJSONBody(t *testing.T, v any) io.Reader {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}
	return bytes.NewReader(b)
}

func TestAccountGet_OK(t *testing.T) {
	t.Parallel()

	uc := &usecaseMock{
		SelectUserByIDFunc: func(userID int) (*entities.User, error) {
			return &entities.User{ID: userID, Email: "x@y.z"}, nil
		},
	}
	s := NewServer("127.0.0.1", 0, uc, "secret", "http://localhost", "./materials")

	req := httptest.NewRequest(http.MethodGet, "/account/my", nil)
	req.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 42))
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestAccountPasswordUpdate_BadMultipart(t *testing.T) {
	t.Parallel()

	s := NewServer("127.0.0.1", 0, &usecaseMock{}, "secret", "http://localhost", "./materials")
	req := httptest.NewRequest(http.MethodPut, "/account/password", bytes.NewBufferString("x"))
	req.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestAccountPasswordUpdate_PasswordsDoNotMatch(t *testing.T) {
	t.Parallel()

	uc := &usecaseMock{}
	s := NewServer("127.0.0.1", 0, uc, "secret", "http://localhost", "./materials")

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("old_password", "old")
	_ = w.WriteField("new_password", "new1")
	_ = w.WriteField("new_password_confirm", "new2")
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPut, "/account/password", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestAccountPasswordUpdate_OK(t *testing.T) {
	t.Parallel()

	uc := &usecaseMock{
		SelectUserByIDFunc: func(userID int) (*entities.User, error) {
			return &entities.User{ID: userID, Email: "a@b.c"}, nil
		},
		AuthorizeFunc: func(user entities.User) (*bool, error) {
			// handler подставляет old_password в user.Password
			if user.Password != "old" {
				return boolPtr(false), nil
			}
			return boolPtr(true), nil
		},
		UpdateUserPasswordFunc: func(userID int, newPassword string) error {
			if userID != 1 || newPassword != "new" {
				return errors.New("unexpected args")
			}
			return nil
		},
	}
	s := NewServer("127.0.0.1", 0, uc, "secret", "http://localhost", "./materials")

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("old_password", "old")
	_ = w.WriteField("new_password", "new")
	_ = w.WriteField("new_password_confirm", "new")
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPut, "/account/password", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestActivityLogList_InvalidUserIDQuery(t *testing.T) {
	t.Parallel()

	s := NewServer("127.0.0.1", 0, &usecaseMock{}, "secret", "http://localhost", "./materials")
	req := httptest.NewRequest(http.MethodGet, "/activitylog?user_id=abc", nil)
	req.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestUsersList_OK_AndErr(t *testing.T) {
	t.Parallel()

	s1 := NewServer("127.0.0.1", 0, &usecaseMock{
		GetUsersFunc: func() ([]entities.User, error) { return []entities.User{{ID: 1}}, nil },
	}, "secret", "http://localhost", "./materials")

	req := httptest.NewRequest(http.MethodGet, "/user/list", nil)
	req.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	rec := httptest.NewRecorder()
	s1.server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	s2 := NewServer("127.0.0.1", 0, &usecaseMock{
		GetUsersFunc: func() ([]entities.User, error) { return nil, errors.New("db") },
	}, "secret", "http://localhost", "./materials")
	rec2 := httptest.NewRecorder()
	s2.server.ServeHTTP(rec2, req)
	if rec2.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d body=%s", rec2.Code, rec2.Body.String())
	}
}

func TestUserCreate_ValidationBranches(t *testing.T) {
	t.Parallel()

	uc := &usecaseMock{
		CreateUserFunc: func(user entities.User) (*entities.User, error) {
			user.ID = 1
			return &user, nil
		},
	}
	s := NewServer("127.0.0.1", 0, uc, "secret", "http://localhost", "./materials")

	// invalid json => 400
	reqBad := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString("{"))
	reqBad.Header.Set("Content-Type", "application/json")
	reqBad.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	recBad := httptest.NewRecorder()
	s.server.ServeHTTP(recBad, reqBad)
	if recBad.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", recBad.Code, recBad.Body.String())
	}

	// missing phone => 422
	req422 := httptest.NewRequest(http.MethodPost, "/user", mustJSONBody(t, map[string]any{
		"email": "a@b.c", "password": "Test12345", "last_name": "L", "first_name": "F", "role": "Админ",
	}))
	req422.Header.Set("Content-Type", "application/json")
	req422.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	rec422 := httptest.NewRecorder()
	s.server.ServeHTTP(rec422, req422)
	if rec422.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d body=%s", rec422.Code, rec422.Body.String())
	}

	// ok
	reqOK := httptest.NewRequest(http.MethodPost, "/user", mustJSONBody(t, map[string]any{
		"email": "a@b.c", "password": "Test12345", "last_name": "L", "first_name": "F", "role": "Админ", "phone": "1",
	}))
	reqOK.Header.Set("Content-Type", "application/json")
	reqOK.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	recOK := httptest.NewRecorder()
	s.server.ServeHTTP(recOK, reqOK)
	if recOK.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recOK.Code, recOK.Body.String())
	}
}

func TestUserActive_AndUserRoleUpdate_QueryValidation(t *testing.T) {
	t.Parallel()

	s := NewServer("127.0.0.1", 0, &usecaseMock{}, "secret", "http://localhost", "./materials")
	tok := signTestToken(t, "secret", 1)

	req := httptest.NewRequest(http.MethodPut, "/user/active", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}

	req2 := httptest.NewRequest(http.MethodPut, "/user/role?email=a@b.c", nil)
	req2.Header.Set("Authorization", "Bearer "+tok)
	rec2 := httptest.NewRecorder()
	s.server.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec2.Code, rec2.Body.String())
	}
}

func TestStatusGet_OK(t *testing.T) {
	t.Parallel()

	s := NewServer("127.0.0.1", 0, &usecaseMock{
		GetStatusesFunc: func() ([]entities.Status, error) { return []entities.Status{{ID: 1}}, nil },
	}, "secret", "http://localhost", "./materials")

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	req.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestFindUser_BindErr(t *testing.T) {
	t.Parallel()

	s := NewServer("127.0.0.1", 0, &usecaseMock{}, "secret", "http://localhost", "./materials")
	req := httptest.NewRequest(http.MethodPost, "/finduser", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestTaskHandlers_StatusParamRequired(t *testing.T) {
	t.Parallel()

	s := NewServer("127.0.0.1", 0, &usecaseMock{}, "secret", "http://localhost", "./materials")
	req := httptest.NewRequest(http.MethodPut, "/task/1/status", nil)
	req.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestEditStatusUpdate_StatusParamRequired(t *testing.T) {
	t.Parallel()

	s := NewServer("127.0.0.1", 0, &usecaseMock{}, "secret", "http://localhost", "./materials")
	req := httptest.NewRequest(http.MethodPut, "/edit/1/status", nil)
	req.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestTaskHandlers_FullFlow(t *testing.T) {
	t.Parallel()

	uc := &usecaseMock{
		CreateTaskFunc: func(task entities.Tasks) (*entities.Tasks, error) {
			task.ID = 1
			return &task, nil
		},
		TaskDeleteFunc: func(id int) error {
			if id == 2 {
				return errors.New("cannot delete")
			}
			return nil
		},
		TaskGetByIdFunc: func(taskId int) (*entities.Tasks, error) {
			if taskId == 3 {
				return nil, errors.New("not found")
			}
			return &entities.Tasks{ID: taskId}, nil
		},
		TaskStatusUpdateFunc: func(taskId int, status string) error {
			if status == "ERR" {
				return errors.New("bad")
			}
			return nil
		},
		TasksListFunc: func(userID int) ([]entities.Tasks, error) {
			return []entities.Tasks{{ID: 10}}, nil
		},
	}
	s := NewServer("127.0.0.1", 0, uc, "secret", "http://localhost", "./materials")
	tok := signTestToken(t, "secret", 1)

	// create bad json
	reqBad := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBufferString("{"))
	reqBad.Header.Set("Content-Type", "application/json")
	reqBad.Header.Set("Authorization", "Bearer "+tok)
	recBad := httptest.NewRecorder()
	s.server.ServeHTTP(recBad, reqBad)
	if recBad.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", recBad.Code, recBad.Body.String())
	}

	// create ok
	reqOK := httptest.NewRequest(http.MethodPost, "/task", mustJSONBody(t, map[string]any{
		"title": "t", "description": "d", "priority": 1, "id_status": 1,
	}))
	reqOK.Header.Set("Content-Type", "application/json")
	reqOK.Header.Set("Authorization", "Bearer "+tok)
	recOK := httptest.NewRecorder()
	s.server.ServeHTTP(recOK, reqOK)
	if recOK.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recOK.Code, recOK.Body.String())
	}

	// delete bad id
	reqDelBad := httptest.NewRequest(http.MethodDelete, "/task/abc", nil)
	reqDelBad.Header.Set("Authorization", "Bearer "+tok)
	recDelBad := httptest.NewRecorder()
	s.server.ServeHTTP(recDelBad, reqDelBad)
	if recDelBad.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", recDelBad.Code, recDelBad.Body.String())
	}

	// delete error branch
	reqDelErr := httptest.NewRequest(http.MethodDelete, "/task/2", nil)
	reqDelErr.Header.Set("Authorization", "Bearer "+tok)
	recDelErr := httptest.NewRecorder()
	s.server.ServeHTTP(recDelErr, reqDelErr)
	if recDelErr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", recDelErr.Code, recDelErr.Body.String())
	}

	// get bad id
	reqGetBad := httptest.NewRequest(http.MethodGet, "/task/abc", nil)
	reqGetBad.Header.Set("Authorization", "Bearer "+tok)
	recGetBad := httptest.NewRecorder()
	s.server.ServeHTTP(recGetBad, reqGetBad)
	if recGetBad.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", recGetBad.Code, recGetBad.Body.String())
	}

	// get error from usecase
	reqGetErr := httptest.NewRequest(http.MethodGet, "/task/3", nil)
	reqGetErr.Header.Set("Authorization", "Bearer "+tok)
	recGetErr := httptest.NewRecorder()
	s.server.ServeHTTP(recGetErr, reqGetErr)
	if recGetErr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", recGetErr.Code, recGetErr.Body.String())
	}

	// status update ok
	reqStatus := httptest.NewRequest(http.MethodPut, "/task/1/status", bytes.NewBufferString("status=OK"))
	reqStatus.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqStatus.Header.Set("Authorization", "Bearer "+tok)
	recStatus := httptest.NewRecorder()
	s.server.ServeHTTP(recStatus, reqStatus)
	if recStatus.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recStatus.Code, recStatus.Body.String())
	}

	// status update usecase error
	reqStatus2 := httptest.NewRequest(http.MethodPut, "/task/1/status", bytes.NewBufferString("status=ERR"))
	reqStatus2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqStatus2.Header.Set("Authorization", "Bearer "+tok)
	recStatus2 := httptest.NewRecorder()
	s.server.ServeHTTP(recStatus2, reqStatus2)
	if recStatus2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", recStatus2.Code, recStatus2.Body.String())
	}

	// list
	reqList := httptest.NewRequest(http.MethodGet, "/task/list", nil)
	reqList.Header.Set("Authorization", "Bearer "+tok)
	recList := httptest.NewRecorder()
	s.server.ServeHTTP(recList, reqList)
	if recList.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recList.Code, recList.Body.String())
	}
}

func TestVersionAndRevisionHandlers_OKAndErrors(t *testing.T) {
	t.Parallel()

	uc := &usecaseMock{
		VersionTaskFunc: func(version entities.Version) (*entities.Version, error) { version.ID = 1; return &version, nil },
		VersionsListFunc: func(taskId int) ([]entities.Version, error) {
			if taskId == 2 {
				return nil, errors.New("bad")
			}
			return []entities.Version{{ID: 1}}, nil
		},
		VersionByIdFunc: func(versionId int) (*entities.Version, error) { return &entities.Version{ID: versionId}, nil },

		CreateRevisionFun: func(revision entities.Revision) (*entities.Revision, error) { revision.ID = 1; return &revision, nil },
		EditStatusUpdateF: func(editId int, status string) error { return nil },
		GetRevisionsByVer: func(versionID int) ([]entities.Revision, error) { return []entities.Revision{{ID: 1}}, nil },
		GetRevisionByIDF:  func(revisionID int) (*entities.Revision, error) { return &entities.Revision{ID: revisionID}, nil },
	}
	s := NewServer("127.0.0.1", 0, uc, "secret", "http://localhost", "./materials")
	tok := signTestToken(t, "secret", 1)

	// version create bad id
	reqBad := httptest.NewRequest(http.MethodPost, "/version/abc", mustJSONBody(t, map[string]any{"title": "t", "description": "d"}))
	reqBad.Header.Set("Content-Type", "application/json")
	reqBad.Header.Set("Authorization", "Bearer "+tok)
	recBad := httptest.NewRecorder()
	s.server.ServeHTTP(recBad, reqBad)
	if recBad.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", recBad.Code, recBad.Body.String())
	}

	// version create ok
	reqOK := httptest.NewRequest(http.MethodPost, "/version/1", mustJSONBody(t, map[string]any{"title": "t", "description": "d"}))
	reqOK.Header.Set("Content-Type", "application/json")
	reqOK.Header.Set("Authorization", "Bearer "+tok)
	recOK := httptest.NewRecorder()
	s.server.ServeHTTP(recOK, reqOK)
	if recOK.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recOK.Code, recOK.Body.String())
	}

	// versions list error
	reqListErr := httptest.NewRequest(http.MethodGet, "/version/list/2", nil)
	reqListErr.Header.Set("Authorization", "Bearer "+tok)
	recListErr := httptest.NewRecorder()
	s.server.ServeHTTP(recListErr, reqListErr)
	if recListErr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", recListErr.Code, recListErr.Body.String())
	}

	// version get ok
	reqGet := httptest.NewRequest(http.MethodGet, "/version/1", nil)
	reqGet.Header.Set("Authorization", "Bearer "+tok)
	recGet := httptest.NewRecorder()
	s.server.ServeHTTP(recGet, reqGet)
	if recGet.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recGet.Code, recGet.Body.String())
	}

	// edit create ok
	reqEdit := httptest.NewRequest(http.MethodPost, "/edit/1", mustJSONBody(t, map[string]any{"title": "t", "description": "d"}))
	reqEdit.Header.Set("Content-Type", "application/json")
	reqEdit.Header.Set("Authorization", "Bearer "+tok)
	recEdit := httptest.NewRecorder()
	s.server.ServeHTTP(recEdit, reqEdit)
	if recEdit.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recEdit.Code, recEdit.Body.String())
	}

	// edit list ok
	reqEdits := httptest.NewRequest(http.MethodGet, "/edit/list/1", nil)
	reqEdits.Header.Set("Authorization", "Bearer "+tok)
	recEdits := httptest.NewRecorder()
	s.server.ServeHTTP(recEdits, reqEdits)
	if recEdits.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recEdits.Code, recEdits.Body.String())
	}
}

func TestMaterialDownload_FileNotFound(t *testing.T) {
	// Не параллелим: меняем cwd.
	chdirMu.Lock()
	defer chdirMu.Unlock()

	tmp := t.TempDir()
	old, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(old) }()

	s := NewServer("127.0.0.1", 0, &usecaseMock{
		GetMaterialFunc: func(materialId int) (*entities.Material, error) {
			return &entities.Material{ID: materialId, Title: "x", Extension: "txt"}, nil
		},
	}, "secret", "http://localhost", "./materials")

	req := httptest.NewRequest(http.MethodGet, "/material/download/1", nil)
	req.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestMaterialDownload_OK(t *testing.T) {
	// Не параллелим: меняем cwd и пишем файлы.
	chdirMu.Lock()
	defer chdirMu.Unlock()

	tmp := t.TempDir()
	old, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(old) }()

	_ = os.MkdirAll("materials", 0o755)
	if err := os.WriteFile(filepath.Join("materials", "file-1.txt"), []byte("ok"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	s := NewServer("127.0.0.1", 0, &usecaseMock{
		GetMaterialFunc: func(materialId int) (*entities.Material, error) {
			return &entities.Material{ID: materialId, Title: "file", Extension: "txt"}, nil
		},
	}, "secret", "http://localhost", "./materials")

	req := httptest.NewRequest(http.MethodGet, "/material/download/1", nil)
	req.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 1))
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("Content-Disposition") == "" {
		t.Fatalf("expected Content-Disposition header")
	}
}

func TestMaterialUpload_Get_List_OK(t *testing.T) {
	// Не параллелим: меняем cwd и пишем файлы.
	chdirMu.Lock()
	defer chdirMu.Unlock()

	tmp := t.TempDir()
	old, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(old) }()

	uc := &usecaseMock{
		CreateMaterialFunc: func(material entities.Material) (int, error) {
			if material.TaskID != 1 || material.CreatorID != 1 || material.Title == "" || material.Extension == "" {
				return 0, errors.New("bad material")
			}
			return 5, nil
		},
		GetMaterialFunc: func(materialId int) (*entities.Material, error) {
			return &entities.Material{ID: materialId, Title: "file", Extension: "txt"}, nil
		},
		GetMaterialsByTaskIDFun: func(taskID int) ([]entities.Material, error) {
			return []entities.Material{{ID: 1}}, nil
		},
	}
	s := NewServer("127.0.0.1", 0, uc, "secret", "http://localhost", "./materials")
	tok := signTestToken(t, "secret", 1)

	// multipart upload
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("description", "desc")
	fw, err := w.CreateFormFile("file", "file.txt")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	_, _ = fw.Write([]byte("content"))
	_ = w.Close()

	reqUp := httptest.NewRequest(http.MethodPost, "/material/1", &buf)
	reqUp.Header.Set("Content-Type", w.FormDataContentType())
	reqUp.Header.Set("Authorization", "Bearer "+tok)
	recUp := httptest.NewRecorder()
	s.server.ServeHTTP(recUp, reqUp)
	if recUp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recUp.Code, recUp.Body.String())
	}

	// file saved
	if _, err := os.Stat(filepath.Join("materials", "file-5.txt")); err != nil {
		t.Fatalf("expected saved file, err=%v", err)
	}

	// material get
	reqGet := httptest.NewRequest(http.MethodGet, "/material/5", nil)
	reqGet.Header.Set("Authorization", "Bearer "+tok)
	recGet := httptest.NewRecorder()
	s.server.ServeHTTP(recGet, reqGet)
	if recGet.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recGet.Code, recGet.Body.String())
	}

	// list
	reqList := httptest.NewRequest(http.MethodGet, "/material/list/1", nil)
	reqList.Header.Set("Authorization", "Bearer "+tok)
	recList := httptest.NewRecorder()
	s.server.ServeHTTP(recList, reqList)
	if recList.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recList.Code, recList.Body.String())
	}
}

func TestNotImplemented(t *testing.T) {
	t.Parallel()

	// Быстро проверяем заглушку напрямую, без сервера.
	e := NewServer("127.0.0.1", 0, &usecaseMock{}, "secret", "http://localhost", "./materials").server
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := notImplemented(c); err != nil {
		t.Fatalf("notImplemented: %v", err)
	}
	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501, got %d", rec.Code)
	}
}
