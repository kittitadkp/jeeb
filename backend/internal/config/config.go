package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   ServerConfig
	MongoDB  MongoConfig
	Keycloak KeycloakConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port         string `envconfig:"PORT" default:"8080"`
	ReadTimeout  int    `envconfig:"READ_TIMEOUT" default:"10"`
	WriteTimeout int    `envconfig:"WRITE_TIMEOUT" default:"10"`
}

type MongoConfig struct {
	URI      string `envconfig:"MONGO_URI" required:"true"`
	Database string `envconfig:"MONGO_DATABASE" default:"jeeb"`
}

type KeycloakConfig struct {
	URL      string `envconfig:"KEYCLOAK_URL" required:"true"`
	Realm    string `envconfig:"KEYCLOAK_REALM" required:"true"`
	ClientID string `envconfig:"KEYCLOAK_CLIENT_ID" required:"true"`
}

type LogConfig struct {
	Level string `envconfig:"LOG_LEVEL" default:"DEBUG"`
}

func Load() (*Config, error) {
	loadEnvFile()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func loadEnvFile() {
	env := os.Getenv("GO_ENV")

	var filename string
	if env == "" {
		filename = "env/.env.local"
	} else {
		filename = fmt.Sprintf("env/.env.%s", env)
	}

	if err := godotenv.Load(filename); err != nil {
		slog.Warn("env file not loaded", "file", filename, "reason", err.Error())
	} else {
		slog.Info("loaded env file", "file", filename)
	}
}
