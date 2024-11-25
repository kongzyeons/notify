package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	LineToken       string
	BitkubAPIKey    string
	BitkubAPISecret string
}

func InitConfig() Config {
	return Config{
		LineToken:       getEnvString("LINE_TOKEN", ""),
		BitkubAPIKey:    getEnvString("BITKUB_API_KEY", ""),
		BitkubAPISecret: getEnvString("BITKUB_API_SECRET", ""),
	}
}

//lint:ignore U1000 Ignore unused code, it may require in the future
func getEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

//lint:ignore U1000 Ignore unused code, it may require in the future
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

//lint:ignore U1000 Ignore unused code, it may require in the future
func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}

//lint:ignore U1000 Ignore unused code, it may require in the future
func getEnvStringArray(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	result := strings.Split(value, ",")
	for i := range result {
		result[i] = strings.TrimSpace(result[i])
	}

	return result
}

//lint:ignore U1000 Ignore unused code, it may require in the future
func getEnvDurationFromSeconds(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		return defaultValue
	}

	return time.Duration(intValue) * time.Second
}

//lint:ignore U1000 Ignore unused code, it may require in the future
func getEnvDurationFromSecondsNullable(key string, defaultValue time.Duration) *time.Duration {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue == 0 {
			return nil
		} else {
			return &defaultValue
		}
	}

	intValue, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		return &defaultValue
	}

	result := time.Duration(intValue) * time.Second
	return &result
}
