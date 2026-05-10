package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Создание задачи
func (s *Server) taskCreate(e echo.Context) error {
	var task entities.Tasks

	userID, err := s.userIDFromToken(e)
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
	fmt.Println(err)
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
		return e.JSON(http.StatusBadRequest, map[string]any{"error": "Невозможно удалить задачу когда она взята в работу"})
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
	userID, err := s.userIDFromToken(e)
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
