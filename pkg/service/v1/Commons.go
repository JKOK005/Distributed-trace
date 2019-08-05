package v1

import (
	"os"
	"strconv"
	"strings"
)

func getEnvStr(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	valueInt, _ := strconv.Atoi(value)
	return valueInt
}

func getEnvStrSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return strings.Split(value, ",")
}