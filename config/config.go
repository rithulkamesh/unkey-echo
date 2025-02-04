package config

import "os"

type Config struct {
	MongoURL      string
	DatabaseName  string
	JWTSecret     string
	RedisURL      string
	RedisPassword string
	AllowedHosts  []string
}

func LoadConfig() *Config {
	return &Config{
		MongoURL:      getEnvOrDefault("MONGO_URL", "mongodb://localhost:27017"),
		DatabaseName:  getEnvOrDefault("DB_NAME", "marketplace"),
		JWTSecret:     getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		RedisURL:      getEnvOrDefault("REDIS_URL", "localhost:6379"),
		RedisPassword: getEnvOrDefault("REDIS_PASSWORD", ""),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
