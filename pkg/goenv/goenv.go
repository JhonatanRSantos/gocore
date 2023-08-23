package goenv

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type (
	Env string
)

var (
	Test        Env = "test"
	Local       Env = "local"
	Development Env = "development"
	Staging     Env = "staging"
	Production  Env = "production"
)

// Load Load an environment variable
func Load[T int64 | float64 | string | bool](name string, defaultValue T) T {
	if envValue := os.Getenv(name); envValue != "" {
		var err error
		var value any

		switch any(defaultValue).(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
			value, err = strconv.ParseInt(envValue, 10, 64)
		case float32, float64:
			value, err = strconv.ParseFloat(envValue, 64)
		case bool:
			value = strings.ToLower(envValue) == "true"
		default:
			value = envValue
		}

		if err != nil {
			panic(fmt.Errorf("failed to load the environment variable [%s]. Cause: %s", name, err))
		}
		return any(value).(T)
	}
	return defaultValue
}
