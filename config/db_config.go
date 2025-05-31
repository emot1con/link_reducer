package config

import "os"

type Config struct {
	PublicHost            string
	Port                  string
	DBUser                string
	DBPassword            string
	DBName                string
	JWTExpirationInSecond int64
	JWTSecret             string
	DatabaseURL           string
	AppEnvironment        string
}

func InitConfig() *Config {
	return &Config{
		PublicHost:     GetEnv("PGHOST", "localhost"),
		Port:           GetEnv("PGPORT", "5432"),
		DBUser:         GetEnv("PGUSER", "psql"),
		DBPassword:     GetEnv("POSTGRES_PASSWORD", "psql"),
		DBName:         GetEnv("POSTGRES_DB", "golang"),
		JWTSecret:      GetEnv("JWT_SECRET_KEY", ""),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		AppEnvironment: GetEnv("APP_ENV", ""),
	}
}

func GetEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}
