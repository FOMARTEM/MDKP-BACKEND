package config

// Config - основная структура конфигурации приложения
type Config struct {
	IP   string `yaml:"ip"`   // IP-адрес сервера
	Port int    `yaml:"port"` // Порт сервера

	API     APIConfig     `yaml:"api"`     // Конфигурация API
	Usecase UsecaseConfig `yaml:"usecase"` // Конфигурация бизнес-логики
	DB      DBConfig      `yaml:"db"`      // Конфигурация базы данных
}

// APIConfig - конфигурация для API-слоя
type APIConfig struct {
	SecretKey             string `yaml:"secret_key"`               // Секретный ключ для JWT
	AccessTokenTTLMinutes int    `yaml:"access_token_ttl_minutes"` // Время жизни access-токена в минутах
	RefreshTokenTTLDays   int    `yaml:"refresh_token_ttl_days"`   // Время жизни refresh-токена в днях
}

// UsecaseConfig - конфигурация для слоя бизнес-логики
type UsecaseConfig struct {
	DefaultMessage string `yaml:"default_message"` // Сообщение по умолчанию для health-check
}

// DBConfig - конфигурация для подключения к базе данных
type DBConfig struct {
	Host     string `yaml:"host"`     // Хост базы данных
	Port     int    `yaml:"port"`     // Порт базы данных
	User     string `yaml:"user"`     // Пользователь базы данных
	Password string `yaml:"password"` // Пароль базы данных
	DBname   string `yaml:"dbname"`   // Имя базы данных
	SSLMode  string `yaml:"sslmode"`  // Режим SSL для PostgreSQL
}
