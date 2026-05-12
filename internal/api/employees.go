package api

import (
	"net/http"
	"strconv"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

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

	userID, err := s.userIDFromToken(e)
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
	userID, err := s.userIDFromToken(e)
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
		_, err := s.userIDFromToken(e)
		if err != nil {
			return err
		}
	*/
	limit, offset := s.getLimitOffset(e)

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
	//limit, offset := s.getLimitOffset(e)

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

	if user.Phone == "" {
		return e.JSON(http.StatusUnprocessableEntity, map[string]any{"error": "Введите телефон"})
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
