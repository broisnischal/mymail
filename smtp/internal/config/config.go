package config

import (
	"os"
	"strconv"
)

type Config struct {
	SMTP      SMTPConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	MinIO     MinIOConfig
	TLS       TLSConfig
	DKIM      DKIMConfig
	RateLimit RateLimitConfig
	TempMail  TempMailConfig
}

type SMTPConfig struct {
	Host           string
	Port           int
	Domain         string
	MaxMessageSize int64
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}

type DKIMConfig struct {
	Enabled    bool
	PrivateKey string
	Selector   string
	Domain     string
}

type RateLimitConfig struct {
	EmailsPerUser    int
	EmailsPerHour    int
	ConnectionsPerIP int
}

type TempMailConfig struct {
	Enabled bool
	TTL     int
}

func Load(configPath string) *Config {
	return &Config{
		SMTP: SMTPConfig{
			Host:           getEnv("SMTP_HOST", "0.0.0.0"),
			Port:           getEnvInt("SMTP_PORT", 25),
			Domain:         getEnv("SMTP_DOMAIN", "mymail.com"),
			MaxMessageSize: int64(getEnvInt("SMTP_MAX_SIZE", 10485760)),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgresql://postgres:postgres@postgres:5432/mymail?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://redis:6379"),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "minio:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "mails"),
			UseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
		},
		TLS: TLSConfig{
			Enabled:  getEnv("TLS_ENABLED", "false") == "true",
			CertFile: getEnv("TLS_CERT_FILE", ""),
			KeyFile:  getEnv("TLS_KEY_FILE", ""),
		},
		DKIM: DKIMConfig{
			Enabled:    getEnv("DKIM_ENABLED", "false") == "true",
			PrivateKey: getEnv("DKIM_PRIVATE_KEY", ""),
			Selector:   getEnv("DKIM_SELECTOR", "default"),
			Domain:     getEnv("DKIM_DOMAIN", getEnv("SMTP_DOMAIN", "mymail.com")),
		},
		RateLimit: RateLimitConfig{
			EmailsPerUser:    getEnvInt("RATE_LIMIT_EMAILS_PER_USER", 1000),
			EmailsPerHour:    getEnvInt("RATE_LIMIT_EMAILS_PER_HOUR", 100),
			ConnectionsPerIP: getEnvInt("RATE_LIMIT_CONNECTIONS_PER_IP", 10),
		},
		TempMail: TempMailConfig{
			Enabled: getEnv("TEMP_MAIL_ENABLED", "true") != "false",
			TTL:     getEnvInt("TEMP_MAIL_TTL", 86400),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
