package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres PostgresConfig
	Schedule ScheduleConfig
	Server   ServerConfig
}

type ScheduleConfig struct {
	Time string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type ServerConfig struct {
	Host      string
	Port      string
	DBTimeout time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnv("POSTGRES_DB", "currencies"),
		},
		Schedule: ScheduleConfig{
			Time: getEnv("SCHEDULE_TIME", "0 59 7,15,23 * * *"),
		},
		Server: ServerConfig{
			Host:      getEnv("SERVER_HOST", "0.0.0.0"),
			Port:      getEnv("SERVER_PORT", "8080"),
			DBTimeout: time.Second * 5,
		},
	}
}

func (p PostgresConfig) ConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.DBName)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
