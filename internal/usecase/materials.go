package usecase

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

func (u *Usecase) CreateMaterial(material entities.Material) (int, error) {
	createdMaterial, err := u.p.CreateMaterial(material)

	if err != nil {
		return 0, err
	}

	return createdMaterial.ID, nil
}

func (u *Usecase) GetMaterial(materialId int) (*entities.Material, error) {
	material, err := u.p.GetMaterialByID(materialId)

	if err != nil {
		return nil, err
	}

	return material, nil
}

func (u *Usecase) GetMaterialsByTaskID(taskID int) ([]entities.Material, error) {
	return u.p.GetMaterialsByTaskID(taskID)
}
