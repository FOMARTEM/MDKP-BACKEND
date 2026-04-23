package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadConfig - загружает конфигурацию из YAML-файла
func LoadConfig(pathToFile string) (*Config, error) {
	// Получаем абсолютный путь к файлу
	filename, err := filepath.Abs(pathToFile)
	if err != nil {
		return nil, err
	}

	// Читаем содержимое файла
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config

	// Разбираем YAML в структуру Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
