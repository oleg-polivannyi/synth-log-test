package config

// Read values from ENV and provide default values for missing ones
import (
	"os"
	"strconv"
)

type Config struct {
	Env            string
	Tag            string
	TargetURL      string
	EventFrequency int
	FileName       string
	StdOut         bool
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func LoadConfig() Config {
	return Config{
		Env:            getEnv("ENV", "dev"),
		Tag:            getEnv("TAG", "default-tag"),
		TargetURL:      getEnv("TARGET_URL", "http://localhost:8080"),
		EventFrequency: getEnvInt("EVENT_FREQUENCY", 60),
		FileName:       getEnv("FILE_NAME", ""),
		StdOut:         getEnvBool("STD_OUT", true),
	}
}
