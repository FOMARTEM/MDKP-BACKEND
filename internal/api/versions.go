package api

import (
	"net/http"
	"strconv"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	"github.com/labstack/echo/v4"
)

// Создание версии к задаче
func (s *Server) versionCreate(e echo.Context) error {
	// Данные по версии
	// Ложим в бд
	// Получаем по ней все данные
	var version entities.Version

	userID, err := s.userIDFromToken(e)
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
		userID, err := s.userIDFromToken(e)
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
