package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (s *Server) healthCheck(e echo.Context) error {
	return e.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

func notImplemented(e echo.Context) error {
	return e.JSON(http.StatusNotImplemented, map[string]any{
		"error":  "not implemented",
		"method": e.Request().Method,
		"path":   e.Path(),
	})
}

// Аутентификация и получение токена JWT ++
func (s *Server) authLogin(e echo.Context) error {
	var user entities.User

	err := e.Bind(&user)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	login, err := s.uc.Authorize(user)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, err.Error())
	}

	if !*login {
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	u, err := s.uc.SelectUserByEmail(user.Email)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	claims := jwtCustomClaims{
		UserID: u.ID,
		Role:   u.Role,
		Email:  u.Email,
		Name:   fmt.Sprintf("%s %s", u.LastName, u.FirstName),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8000)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	u.Token, err = token.SignedString([]byte(s.secretKey))

	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, u)
}

//func (s *Server) authRefresh(e echo.Context) error { return notImplemented(e) }
//func (s *Server) authLogout(e echo.Context) error  { return notImplemented(e) }

// Обновление пароля пользователя ++
func (s *Server) accountPasswordUpdate(e echo.Context) error {
	var oldPassword, newPassword, newPasswordConfirm string

	// ПОлучаем данные из multipart/form-data
	form, err := e.MultipartForm()
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	oldPassword = form.Value["old_password"][0]
	newPassword = form.Value["new_password"][0]
	newPasswordConfirm = form.Value["new_password_confirm"][0]

	if newPassword != newPasswordConfirm {
		return e.JSON(http.StatusBadRequest, entities.ErrPasswordsDoNotMatch)
	}

	userID, err := userIDFromToken(e)
	if err != nil {
		return err
	}

	user, err := s.uc.SelectUserByID(userID)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	user.Password = oldPassword

	login, err := s.uc.Authorize(*user)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, err.Error())
	}

	if !*login {
		return e.JSON(http.StatusUnauthorized, entities.ErrInvalidPassword)
	}

	user.Password = newPassword

	err = s.uc.UpdateUserPassword(user.ID, user.Password)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// Получение данных аккаунта текущего пользователя ++
func (s *Server) accountGet(e echo.Context) error {
	userID, err := userIDFromToken(e)
	if err != nil {
		return err
	}

	user, err := s.uc.SelectUserByID(userID)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, user)
}

// Получение Логов, с возможностью фильтрации по дате, пользователю и т.д. ++
func (s *Server) activityLogList(e echo.Context) error {
	// Получаем ID пользователя из контекста, установленного JWT middleware
	/*
		_, err := userIDFromToken(e)
		if err != nil {
			return err
		}
	*/
	limit, offset := getLimitOffset(e)

	userEmail := e.QueryParam("email")
	userIdStr := e.QueryParam("user_id")
	var userId int
	var err error
	if userIdStr != "" {
		userId, err = strconv.Atoi(userIdStr)
		if err != nil {
			return e.JSON(http.StatusBadRequest, entities.ErrInvalidQueryParams)
		}
	}
	startDate := e.QueryParam("start_date")
	endDate := e.QueryParam("end_date")

	var logs []entities.Log

	logs, err = s.uc.GetLogs(userId, userEmail, startDate, endDate, limit, offset)

	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, logs)
}

func (s *Server) usersList(e echo.Context) error {
	//limit, offset := getLimitOffset(e)

	users, err := s.uc.GetUsers()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, users)
}

// Создание сотрудника в системе ++
func (s *Server) userCreate(e echo.Context) error {
	var user entities.User

	err := e.Bind(&user)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	err = validator.New().Struct(user)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	createdUser, err := s.uc.CreateUser(user)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, createdUser)
}

// Изменение фалага активности ++
func (s *Server) userActive(e echo.Context) error {
	userEmail := e.QueryParam("email")

	if userEmail == "" {
		return e.JSON(http.StatusBadRequest, entities.ErrInvalidQueryParams)
	}

	err := s.uc.UserActiveChange(userEmail)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err)
	}

	return e.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

func (s *Server) userRoleUpdate(e echo.Context) error {
	userEmail := e.QueryParam("email")
	roleIDSTR := e.QueryParam("roleID")

	if userEmail == "" || roleIDSTR == "" {
		return e.JSON(http.StatusBadRequest, entities.ErrInvalidQueryParams)
	}

	roleID, err := strconv.Atoi(roleIDSTR)
	if err != nil {
		return e.JSON(http.StatusBadRequest, entities.ErrInvalidQueryParams)
	}

	err = s.uc.UserRoleUpdate(userEmail, roleID)
	if err != nil {
		e.JSON(http.StatusInternalServerError, err)
	}

	return e.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// Получаем список ролей для установки, либо поиска по ним ++
func (s *Server) rolesList(e echo.Context) error {
	userID := e.Get("id")
	if userID == nil {
		return e.JSON(http.StatusUnauthorized, entities.ErrInvalidToken)
	}

	var roles []entities.Role

	roles, err := s.uc.GetRoles()

	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, roles)
}

// Получение возмможных статусов задачи ++
func (s *Server) statusGet(e echo.Context) error {
	var statuses []entities.Status

	statuses, err := s.uc.GetStatuses()
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, statuses)
}

// Ищем пользователя по его Роли ++
func (s *Server) findUser(e echo.Context) error {
	var user entities.User

	err := e.Bind(&user)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	users, err := s.uc.GetUsersByRole(user)

	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, users)
}

// Создание задачи
func (s *Server) taskCreate(e echo.Context) error {
	var task entities.Tasks

	userID, err := userIDFromToken(e)
	if err != nil {
		return err
	}

	err = e.Bind(&task)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	err = validator.New().Struct(task)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	task.IdCreator = userID

	taskCreated, err := s.uc.CreateTask(task)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, taskCreated)
}

func (s *Server) taskDelete(e echo.Context) error {
	taskId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	err = s.uc.TaskDelete(taskId)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

//func (s *Server) taskAssign(e echo.Context) error { return notImplemented(e) }

func (s *Server) taskGet(e echo.Context) error {
	taskId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	task, err := s.uc.TaskGetById(taskId)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, task)
}

func (s *Server) taskStatusUpdate(e echo.Context) error {
	taskId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	status := e.FormValue("status")

	if status == "" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"error": "status parameter is required",
		})
	}

	err = s.uc.TaskStatusUpdate(taskId, status)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

func (s *Server) tasksList(e echo.Context) error {
	userID, err := userIDFromToken(e)
	if err != nil {
		return err
	}

	tasks, err := s.uc.TasksList(userID)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, tasks)
}

//func (s *Server) tasksSearch(e echo.Context) error { return notImplemented(e) }

// Создать материал в базе данных
// Сохранить материал на сервере используя его название и id
// Переписать что бы можно было по одной ручке сохранять либо к задаче либо к версии
func (s *Server) materialUpload(e echo.Context) error {
	var material entities.Material

	err := e.Bind(&material)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	taskId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	userID, err := userIDFromToken(e)
	if err != nil {
		return err
	}

	file, err := e.FormFile("file")
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	filename := file.Filename
	extension := strings.TrimPrefix(filepath.Ext(filename), ".")
	title := strings.TrimSuffix(filename, filepath.Ext(filename))

	/*
		fmt.Println(filename)
		fmt.Println(extension)
		fmt.Println(title)
	*/

	material.Extension = extension
	material.Title = title
	material.TaskID = taskId
	material.CreatorID = userID

	material.ID, err = s.uc.CreateMaterial(material)

	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	// Сохранение материала на сервере
	err = saveMaterialFile(file, material)

	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

func (s *Server) materialGet(e echo.Context) error {
	materialId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	material, err := s.uc.GetMaterial(materialId)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, material)
}

func (s *Server) materialDownload(e echo.Context) error {
	materialId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	material, err := s.uc.GetMaterial(materialId)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	filePath := fmt.Sprintf("./materials/%s-%d.%s", material.Title, material.ID, material.Extension)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return e.JSON(http.StatusNotFound, "file not found")
	}

	newFileName := material.Title + "." + material.Extension

	return e.Attachment(filePath, newFileName)
}

// func (s *Server) materialDelete(e echo.Context) error { return notImplemented(e) }

// Создание версии к задаче
func (s *Server) versionCreate(e echo.Context) error {
	// Данные по версии
	// Ложим в бд
	// Получаем по ней все данные
	var version entities.Version

	userID, err := userIDFromToken(e)
	if err != nil {
		return err
	}

	taskId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	err = e.Bind(&version)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	version.CreatorID = userID
	version.TaskID = taskId

	versionCreated, err := s.uc.VersionTask(version)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, versionCreated)
}

func (s *Server) versionsList(e echo.Context) error {
	/*
		userID, err := userIDFromToken(e)
		if err != nil {
			return err
		}
	*/

	taskId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	versions, err := s.uc.VersionsList(taskId)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, versions)
}

func (s *Server) versionGet(e echo.Context) error {
	versionId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	version, err := s.uc.VersionById(versionId)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, version)
}

func (s *Server) editCreate(e echo.Context) error { return notImplemented(e) }

func (s *Server) editStatusUpdate(e echo.Context) error { return notImplemented(e) }

func (s *Server) editsList(e echo.Context) error { return notImplemented(e) }

func (s *Server) editDelete(e echo.Context) error { return notImplemented(e) }
