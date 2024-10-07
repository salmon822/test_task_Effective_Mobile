package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Server   *ServerConfig
		Postgres *PostgresConfig
		Handler  *HandlerConfig
	}
	PostgresConfig struct {
		Host     string
		User     string
		Password string
		DBName   string
		Port     int
	}
	ServerConfig struct {
		Port           int
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
		MaxHeaderBytes int
	}
	HandlerConfig struct {
		RequestTimeout time.Duration
		QueueSize      int
	}
)

func Init(configPath string) (*Config, error) {
	jsonCfg := viper.New()
	jsonCfg.SetConfigFile(configPath)

	if err := jsonCfg.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config/Init/jsonCfg.ReadInConfig: %w", err)
	}

	envCfg := viper.New()
	envCfg.SetConfigFile(".env")

	if err := envCfg.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config/Init/envCfg.ReadInConfig: %w", err)
	}
	return &Config{
		Server: &ServerConfig{
			Port:           jsonCfg.GetInt("server.port"),
			ReadTimeout:    jsonCfg.GetDuration("server.readTimeout"),
			WriteTimeout:   jsonCfg.GetDuration("server.writeTimeout"),
			MaxHeaderBytes: jsonCfg.GetInt("server.maxHeaderBytes"),
		},
		Postgres: &PostgresConfig{
			Host:     envCfg.GetString("POSTGRES_HOST"),
			User:     envCfg.GetString("POSTGRES_USER"),
			Password: envCfg.GetString("POSTGRES_PASSWORD"),
			DBName:   envCfg.GetString("POSTGRES_DB"),
			Port:     envCfg.GetInt("POSTGRES_PORT"),
		},
		Handler: &HandlerConfig{
			RequestTimeout: jsonCfg.GetDuration("handler.requestTimeout"),
			QueueSize:      jsonCfg.GetInt("handler.queueSize"),
		},
	}, nil
}

func (p *PostgresConfig) PgSource() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.DBName)
}
