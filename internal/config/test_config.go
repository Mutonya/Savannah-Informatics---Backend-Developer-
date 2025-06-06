package config

type TestConfig struct {
	Environment   string
	Database      DatabaseConfig
	OAuth         OAuthConfig
	Server        ServerConfig
	SMTP          SMTPConfig
	AfricaTalking AfricaTalkingConfig
}

type DatabaseConfig struct {
	URL           string
	MaxOpenConns  int
	MaxIdleConns  int
	MaxIdleTime   string
	MigrationPath string
}

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	Scopes       []string
}

type ServerConfig struct {
	Port            string
	Timeout         string
	CookieDomain    string
	CookieSecure    bool
	CookieSecret    string
	AccessTokenTTL  string
	RefreshTokenTTL string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type AfricaTalkingConfig struct {
	APIKey    string
	Username  string
	ShortCode string
}

// LoadTestConfig returns a hardcoded configuration used only for tests.
func LoadTestConfig() *TestConfig {
	return &TestConfig{
		Environment: "test",
		Database: DatabaseConfig{
			URL:           "postgres://postgres:postgres@localhost:5432/savanah_test?sslmode=disable",
			MaxOpenConns:  10,
			MaxIdleConns:  5,
			MaxIdleTime:   "5m",
			MigrationPath: "./migrations",
		},
		OAuth: OAuthConfig{
			ClientID:     "test_client_id",
			ClientSecret: "test_client_secret",
			RedirectURL:  "http://localhost:8080/auth/callback",
			AuthURL:      "https://test-oauth-provider.com/auth",
			TokenURL:     "https://test-oauth-provider.com/token",
			UserInfoURL:  "https://test-oauth-provider.com/userinfo",
			Scopes:       []string{"openid", "profile", "email"},
		},
		Server: ServerConfig{
			Port:            "8080",
			Timeout:         "30s",
			CookieDomain:    "localhost",
			CookieSecure:    false,
			CookieSecret:    "test-secret-123",
			AccessTokenTTL:  "15m",
			RefreshTokenTTL: "24h",
		},
		SMTP: SMTPConfig{
			Host:     "localhost",
			Port:     1025,
			Username: "",
			Password: "",
			From:     "test@example.com",
		},
		AfricaTalking: AfricaTalkingConfig{
			APIKey:    "test_api_key",
			Username:  "test_username",
			ShortCode: "TEST",
		},
	}
}
