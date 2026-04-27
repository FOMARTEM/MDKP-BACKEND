package entities

import "errors"

var (
	// Ошибка доступа
	ErrAccessDenied = errors.New("access denied")
	// Ошибка авторизации, не верный пароль пользователя
	ErrInvalidPassword = errors.New("invalid password")
	// Ошибка авторизации, пользователь не найден
	ErrUserNotFound = errors.New("user not found")
	// Ошибка пароль слишком длинный
	ErrPasswordTooLong = errors.New("password too long")
	// Ошибка пользователь с таким email уже существует
	ErrUserAlreadyExists = errors.New("user with this email already exists")
	// Ошибка пользователь не активен
	ErrUserInactive = errors.New("user is inactive")
	// Ошибка пользователь не авторизован
	ErrUserUnauthorized = errors.New("user is unauthorized")
	// Ошибка токена, не верный токен
	ErrInvalidToken = errors.New("invalid token")
	// Ошибка пароли не совпадают
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
	// Ошибка не верно указаны query параметры
	ErrInvalidQueryParams = errors.New("invalid query parameters")
)
