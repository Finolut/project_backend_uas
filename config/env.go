package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	DBDriver    string
	PostgresDsn string
	MongoURI    string
	JWTSecret   string
	LogPath     string
	LogLevel    string
}

// singleton config
var (
	cfg  *Config
	once sync.Once
)

// LoadEnv reads .env (if present) and environment variables.
// Call this once at startup (or call Get() which ensures LoadEnv ran).
func LoadEnv() error {
	var loadErr error
	once.Do(func() {
		// try load .env silently
		_ = godotenv.Load()

		c := &Config{
			AppPort:     getEnv("APP_PORT", "3000"),
			DBDriver:    getEnv("DB_DRIVER", "mongo"), // mongo or postgres
			PostgresDsn: getEnv("POSTGRES_DSN", ""),
			MongoURI:    getEnv("MONGO_URI", ""),
			JWTSecret:   getEnv("JWT_SECRET", "dev-secret"),
			LogPath:     getEnv("LOG_PATH", "logs/app.log"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		}
		cfg = c
	})
	return loadErr
}

// Get returns loaded config. It ensures LoadEnv was called.
func Get() *Config {
	if cfg == nil {
		if err := LoadEnv(); err != nil {
			log.Printf("warning: LoadEnv returned error: %v", err)
		}
	}
	return cfg
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
