package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Schedule ScheduleConfig
	Server   ServerConfig
}

type ScheduleConfig struct {
	Time string
}

type ServerConfig struct {
	Host string
	Port string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Schedule: ScheduleConfig{
			Time: getEnv("SCHEDULE_TIME", "0 0 * * * *"),
		},
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
