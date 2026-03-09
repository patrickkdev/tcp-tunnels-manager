package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/patrickkdev/tcptunnel/internal/infrastructure/db"
)

var (
	DBConfig                 db.Config
	ReconcileIntervalSeconds = 5
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	DBConfig = db.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "tcp_tunnels"),
		Port:     getEnv("DB_PORT", 3306),
	}

	ReconcileIntervalSeconds = getEnv("RECONCILE_INTERVAL_SECONDS", ReconcileIntervalSeconds)
}

func getEnv[T any](key string, defaultValue T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	var result T
	switch any(result).(type) {
	case int:
		fmt.Sscanf(value, "%d", &result)
	case bool:
		fmt.Sscanf(value, "%t", &result)
	case float64, float32:
		fmt.Sscanf(value, "%f", &result)
	default:
		result = any(value).(T)
	}

	return result
}
