package config

import "os"

type Config struct {
	ListenAddr   string
	DatabasePath string
	CloudAPIUrl  string
}

func Load() *Config {
	return &Config{
		ListenAddr:   getEnv("PRYX_LISTEN_ADDR", ":3000"),
		DatabasePath: getEnv("PRYX_DB_PATH", "pryx.db"),
		CloudAPIUrl:  getEnv("PRYX_CLOUD_API_URL", "https://pryx.dev/api"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
