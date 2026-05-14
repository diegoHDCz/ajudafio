package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	// MigrationsPath string
}

func Load() *Config {
	// err := godotenv.Load()
	// if err != nil {
	// 	panic(fmt.Sprintf("config: failed to load environment variables: %v", err))
	// }
	return &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseURL: mustGetEnv("DATABASE_DOCKER_URL"),
		// MigrationsPath: getEnv("MIGRATIONS_PATH", "./migrations"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		fmt.Printf("config: using environment variable %q=%q\n", key, v)
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	fmt.Printf("config: using environment variable %q=%q\n", key, v)
	if v == "" {
		panic(fmt.Sprintf("config: environment variable %q is required", key))
	}
	return v
}
