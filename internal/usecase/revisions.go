package usecase

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

func (u *Usecase) CreateRevision(revision entities.Revision) (*entities.Revision, error) {
	return u.p.CreateRevision(revision)
}

func (u *Usecase) EditStatusUpdate(editId int, status string) error {
	return u.p.UpdateRevisionStatus(editId, status)
}

func (u *Usecase) GetRevisionsByVersionID(versionID int) ([]entities.Revision, error) {
	return u.p.GetRevisionsByVersionID(versionID)
}

func (u *Usecase) GetRevisionByID(revisionID int) (*entities.Revision, error) {
	return u.p.GetRevisionByID(revisionID)
}
