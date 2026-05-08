package usecase

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

// Работа с версиями

func (u *Usecase) VersionTask(version entities.Version) (*entities.Version, error) {
	versionCreated, err := u.p.CreateVersion(version)

	if err != nil {
		return nil, err
	}

	return versionCreated, err
}

func (u *Usecase) VersionsList(taskId int) ([]entities.Version, error) {
	return u.p.GetVersionsByTaskID(taskId)
}

func (u *Usecase) VersionById(versionId int) (*entities.Version, error) {
	return u.p.GetVersionByID(versionId)
}
