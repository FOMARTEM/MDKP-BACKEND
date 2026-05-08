package api

import (
	"net/http"
	"strconv"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	"github.com/labstack/echo/v4"
)

func (s *Server) editCreate(e echo.Context) error {
	var revision entities.Revision

	userID, err := s.userIDFromToken(e)
	if err != nil {
		return err
	}

	versionId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	err = e.Bind(&revision)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	revision.CreatorID = userID
	revision.VersionID = versionId

	createdRevision, err := s.uc.CreateRevision(revision)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, createdRevision)
}

func (s *Server) editStatusUpdate(e echo.Context) error {
	editId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	status := e.FormValue("status")

	if status == "" {
		return e.JSON(http.StatusBadRequest, map[string]any{
			"error": "status parameter is required",
		})
	}

	err = s.uc.EditStatusUpdate(editId, status)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

func (s *Server) editsList(e echo.Context) error {
	versionId, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	revisions, err := s.uc.GetRevisionsByVersionID(versionId)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, revisions)
}

//func (s *Server) editDelete(e echo.Context) error { return notImplemented(e) }

func (s *Server) editGet(e echo.Context) error {
	revisionID, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	revision, err := s.uc.GetRevisionByID(revisionID)

	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusOK, revision)
}
