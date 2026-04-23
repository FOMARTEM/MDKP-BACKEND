package entities

import "errors"

var (
	// Ошибка доступа
	ErrAccessDenied = errors.New("access denied")
)
