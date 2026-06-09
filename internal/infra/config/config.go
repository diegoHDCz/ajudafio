package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type Config struct {
	AppPort            string
	DatabaseURL        string
	ClerkJWKSURL       string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSS3BucketName    string
	// MigrationsPath string
}

func Load() *Config {
	// if os.Getenv("ENVIRONMENT") == "DEV" {
	_ = godotenv.Load()

	// }
	return &Config{
		AppPort:            getEnv("APP_PORT", "8080"),
		DatabaseURL:        mustGetEnv("DATABASE_URL"),
		ClerkJWKSURL:       mustGetEnv("CLERK_JWKS_URL"),
		AWSAccessKeyID:     mustGetEnv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: mustGetEnv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:          mustGetEnv("AWS_REGION"),
		AWSS3BucketName:    mustGetEnv("AWS_S3_BUCKET_NAME"),
		// MigrationsPath: getEnv("MIGRATIONS_PATH", "./migrations"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		log.Printf("config: using environment variable %q=%q\n", key, v)
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	log.Printf("config: using environment variable %q=%q", key, v)
	if v == "" {
		msg := fmt.Sprintf("config: environment variable %q is required", key)
		log.Print(msg)
		panic(msg)
	}
	return v
}
