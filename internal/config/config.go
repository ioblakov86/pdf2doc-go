package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	MaxUploadSize int64 // bytes
	MaxConcurrent int
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8088"),
		MaxUploadSize: getEnvInt64("MAX_UPLOAD_SIZE", 50<<20), // 50 MB
		MaxConcurrent: getEnvInt("MAX_CONCURRENT", 2),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

func getEnvInt64(key string, def int64) int64 {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return def
}
