package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Schedule ScheduleConfig
}

type ScheduleConfig struct {
	Time string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Schedule: ScheduleConfig{
			Time: getEnv("SCHEDULE_TIME", "* */45 * * * *"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
