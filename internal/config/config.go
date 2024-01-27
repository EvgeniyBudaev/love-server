package config

import (
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

type Config struct {
	Port             string `envconfig:"PORT"`
	LoggerLevel      string `envconfig:"LOGGER_LEVEL"`
	Host             string `envconfig:"HOST"`
	DBPort           string `envconfig:"DB_PORT"`
	DBUser           string `envconfig:"DB_USER"`
	DBPassword       string `envconfig:"DB_PASSWORD"`
	DBName           string `envconfig:"DB_NAME"`
	DBSSlMode        string `envconfig:"DB_SSLMODE"`
	TelegramBotToken string `envconfig:"TELEGRAM_BOT_TOKEN"`
}

func Load(l logger.Logger) (*Config, error) {
	var cfg Config
	if err := godotenv.Load(); err != nil {
		l.Debug("error func Load, method Load by path internal/config/config.go", zap.Error(err))
		return nil, err
	}
	err := envconfig.Process("MYAPP", &cfg)
	if err != nil {
		l.Debug("error func Load, method Process by path internal/config/config.go", zap.Error(err))
		return nil, err
	}
	return &cfg, nil
}
