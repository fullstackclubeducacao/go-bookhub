package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	MongoDB  MongoDBConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Postgres
type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
	Issuer        string
}

type MongoDBConfig struct {
	URI         string
	Database    string
	MaxPoolSize uint64
	MinPoolSize uint64
	MaxIdleTime time.Duration
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getIntEnv("DB_PORT", 5432),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", "postgres"),
			DBName:       getEnv("DB_NAME", "bookhub"),
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns: getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getIntEnv("DB_MAX_IDLE_CONNS", 5),
			MaxLifetime:  getDurationEnv("DB_MAX_LIFETIME", 5*time.Minute),
		},
		MongoDB: MongoDBConfig{
			URI:         getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database:    getEnv("MONGO_DATABASE", "bookhub"),
			MaxPoolSize: getUint64Env("MONGO_MAX_POOL_SIZE", 100),
			MinPoolSize: getUint64Env("MONGO_MIN_POOL_SIZE", 10),
			MaxIdleTime: getDurationEnv("MONGO_MAX_IDLE_TIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			SecretKey:     getEnv("JWT_SECRET_KEY", "your-super-secret-key-change-in-production"),
			TokenDuration: getDurationEnv("JWT_TOKEN_DURATION", 24*time.Hour),
			Issuer:        getEnv("JWT_ISSUER", "bookhub"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getUint64Env(key string, defaultValue uint64) uint64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseUint(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
