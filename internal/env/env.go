package env

import (
	"os"
	"strconv"
	"time"
)

// getEnvInt reads an env var and converts to int, or returns a default
func GetInt(key string, fallback int) int {

	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	intVal, err := strconv.Atoi(os.Getenv(val))
	if err != nil {
		return fallback
	}

	return intVal
}

func GetBool(key string, fallback bool) bool {
	val := os.Getenv(key)

	if val == "" {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolVal
}

// getEnvTime reads an env var and parses it as a duration (e.g. "15m", "1h")
func GetTime(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)

	if val == "" {
		return fallback
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return d
}

func GetString(key, fallback string) string {

	val := os.Getenv(key)

	if val == "" {
		return fallback
	}

	return val
}
