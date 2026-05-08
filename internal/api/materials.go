package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	"github.com/labstack/echo/v4"
)

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

	userID, err := s.userIDFromToken(e)
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
	err = s.saveMaterialFile(file, material)

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

func (s *Server) materialsList(e echo.Context) error {
	taskId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	materials, err := s.uc.GetMaterialsByTaskID(taskId)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, materials)
}

// func (s *Server) materialDelete(e echo.Context) error { return notImplemented(e) }
