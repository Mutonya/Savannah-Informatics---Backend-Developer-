package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string

	OAuthClientID     string
	OAuthClientSecret string
	OAuthRedirectURL  string
	OAuthProviderURL  string

	ServerPort  string
	Environment string

	AfricaTalkingAPIKey   string
	AfricaTalkingUsername string
	SMTPHost              string
	SMTPPort              int
	SMTPUsername          string
	SMTPPassword          string
	AdminEmail            string
	Currency              string
	SMSSenderID           string
}

func LoadConfig() *Config {
	// Get SMTP port as integer
	smtpPortStr := getEnv("SMTP_PORT", "1025") // Default to MailHog port
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		// Fallback to default port if conversion fails
		smtpPort = 1025
	}
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "savannah"),
		SSLMode:    getEnv("SSL_MODE", "disable"),

		OAuthClientID:     getEnv("OAUTH_CLIENT_ID", ""),
		OAuthClientSecret: getEnv("OAUTH_CLIENT_SECRET", ""),
		OAuthRedirectURL:  getEnv("OAUTH_REDIRECT_URL", ""),
		OAuthProviderURL:  getEnv("OAUTH_PROVIDER_URL", ""),

		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		AfricaTalkingAPIKey:   getEnv("AFRICA_TALKING_API_KEY", ""),
		AfricaTalkingUsername: getEnv("AFRICA_TALKING_USERNAME", ""),
		SMSSenderID:           getEnv("CURRENCY", ""),

		SMTPHost: getEnv("SMTP_HOST", ""),

		SMTPPort:     smtpPort,
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		AdminEmail:   getEnv("ADMIN_EMAIL", ""),
		Currency:     getEnv("CURRENCY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
