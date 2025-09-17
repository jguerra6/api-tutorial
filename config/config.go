package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	envFilePath  = ".env"
	configPrefix = "TUTORIAL"
)

type Config struct {
	Port               string   `envconfig:"TUTORIAL_HTTP_SERVER_PORT" required:"true"`
	FirebaseConfigFile string   `envconfig:"TUTORIAL_FIREBASE_CONFIG_FILE" required:"true"`
	AdminTokens        []string `envconfig:"TUTORIAL_AUTH_TOKENS" required:"true"`
	AllowedOrigins     []string `envconfig:"TUTORIAL_ALLOWED_ORIGINS" required:"true"`
	DbUserName         string   `envconfig:"TUTORIAL_DATABASE_USERNAME" required:"true"`
	DbPassword         string   `envconfig:"TUTORIAL_DATABASE_PASSWORD" required:"true"`
	DbHost             string   `envconfig:"TUTORIAL_DATABASE_HOST" required:"true"`
	DbPort             string   `envconfig:"TUTORIAL_DATABASE_PORT" required:"true"`
	DbName             string   `envconfig:"TUTORIAL_DATABASE_NAME" required:"true"`
}

func Load() (*Config, error) {

	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(envFilePath)
		if err != nil {
			return nil, err
		}
	}

	cfg := &Config{}

	err := envconfig.Process(configPrefix, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
