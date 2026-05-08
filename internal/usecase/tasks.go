package usecase

import (
	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
)

// Работа с задачами

// Добавть проверку что пользователь создающий задачу является руководителем, а остальные соотвуствуют другим ролям
// Добавить проверку дедлайна
// Добавить логи
func (u *Usecase) CreateTask(task entities.Tasks) (*entities.Tasks, error) {
	taskCreated, err := u.p.CreateTask(task)

	if err != nil {
		return nil, err
	}

	return taskCreated, nil
}

func (u *Usecase) TaskDelete(id int) error {
	err := u.p.DeleteTask(id)

	return err
}

func (u *Usecase) TaskGetById(taskId int) (*entities.Tasks, error) {
	task, err := u.p.GetTaskByID(taskId)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (u *Usecase) TaskStatusUpdate(taskId int, status string) error {
	err := u.p.UpdateTaskStatus(taskId, status)

	return err
}

func (u *Usecase) TasksList(userID int) ([]entities.Tasks, error) {
	tasks, err := u.p.GetTasksByUserID(userID)

	return tasks, err
}
