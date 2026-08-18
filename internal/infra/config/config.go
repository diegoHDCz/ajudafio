package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type Config struct {
	AppPort              string
	DatabaseURL          string
	JWTSecret            string
	SupabaseAccessKeyID  string
	SupabaseSecretKey    string
	SupabaseRegion       string
	SupabaseBucketName   string
	SupabaseEndpoint     string
	GoogleClientID       string
	// MigrationsPath string
}

func Load() *Config {
	// if os.Getenv("ENVIRONMENT") == "DEV" {
	_ = godotenv.Load()

	// }
	return &Config{
		AppPort:             getEnv("APP_PORT", "8080"),
		DatabaseURL:         mustGetEnv("DATABASE_URL"),
		JWTSecret:           mustGetEnv("JWT_SECRET"),
		SupabaseAccessKeyID: mustGetEnv("SUPABASE_ACCESS_KEY_ID"),
		SupabaseSecretKey:   mustGetEnv("SUPABASE_SECRET_KEY"),
		SupabaseRegion:      mustGetEnv("SUPABASE_REGION"),
		SupabaseBucketName:  mustGetEnv("SUPABASE_BUCKET_NAME"),
		SupabaseEndpoint:    mustGetEnv("SUPABASE_ENDPOINT"),
		GoogleClientID:      getEnv("GOOGLE_CLIENT_ID", ""),
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
