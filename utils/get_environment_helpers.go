package utils

import (
	"os"
	"strconv"
	"strings"
)

// Getenv read an environment variable or return default value
func Getenv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// GetenvAsInt read an environment variable into a int or return default value
// GetenvAsInt("DB_POOL_SIZE", 10)
// DB_POOL_SIZE=20
func GetenvAsInt(name string, defaultVal int) int {
	valueStr := Getenv(name, "")

	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// GetenvAsBool read an environment variable into a bool or return default value
// GetenvAsBool("DEBUG_MODE", false)
// DEBUG_MODE=true
func GetenvAsBool(name string, defaultVal bool) bool {
	valStr := Getenv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// GetenvAsSlice read an environment variable into a string slice or return default value
// GetenvAsSlice("USERS", []string{"admin"}, ",")
// USERS=admin,guest_1,guest_2
func GetenvAsSlice(name string, defaultVal []string, separator string) []string {
	valStr := Getenv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, separator)

	return val
}
