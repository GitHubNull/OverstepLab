package config

import (
	"os"
)

type Config struct {
	Port      string
	DBPath    string
	JWTSecret string
	SafeMode  bool
}

func Load() *Config {
	port := os.Getenv("OVERSTEPLAB_PORT")
	if port == "" {
		port = "5000"
	}

	dbPath := os.Getenv("OVERSTEPLAB_DB_PATH")
	if dbPath == "" {
		dbPath = "oversteplab.db"
	}

	jwtSecret := os.Getenv("OVERSTEPLAB_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "oversteplab-default-jwt-secret-change-in-production"
	}

	safeMode := os.Getenv("OVERSTEPLAB_SAFE_MODE")

	return &Config{
		Port:      port,
		DBPath:    dbPath,
		JWTSecret: jwtSecret,
		SafeMode:  safeMode == "true" || safeMode == "1",
	}
}
