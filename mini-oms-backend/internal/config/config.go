package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	// Server
	Port string
	Env  string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret      string
	JWTExpiryHours int
}

func Load() *Config {
	return &Config{
		// Server
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "mini_oms"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// JWT
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-this"),
		JWTExpiryHours: getEnvAsInt("JWT_EXPIRY_HOURS", 24),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func (c *Config) GetDSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// LogConfig prints configuration (without sensitive data)
func (c *Config) LogConfig() {
	log.Println("Configuration loaded:")
	log.Printf("  Environment: %s", c.Env)
	log.Printf("  Server Port: %s", c.Port)
	log.Printf("  Database: %s@%s:%s/%s", c.DBUser, c.DBHost, c.DBPort, c.DBName)
}
