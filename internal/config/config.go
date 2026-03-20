package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Cookie   CookieConfig
	Log      LogConfig
	CORS     CORSConfig
	SMTP     SMTPConfig
	OAuth    OAuthConfig
	Telegram TelegramConfig
}

type AppConfig struct {
	Name string
	Env  string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	TLSEnabled   bool
	TLSCertFile  string
	TLSKeyFile   string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

type CookieConfig struct {
	Domain   string
	Secure   bool
	SameSite string
	Path     string
}

type LogConfig struct {
	Level  string
	Format string
}

type CORSConfig struct {
	AllowedOrigins []string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type OAuthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

type TelegramConfig struct {
	BotToken string
	ChatID   string
	ProxyURL string
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode,
	)
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = viper.ReadInConfig()

	setDefaults()

	// Validate critical security settings
	if err := validateSecuritySettings(); err != nil {
		return nil, err
	}

	cfg := &Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Env:  viper.GetString("APP_ENV"),
		},
		Server: ServerConfig{
			Port:         viper.GetString("SERVER_PORT"),
			ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
			IdleTimeout:  viper.GetDuration("SERVER_IDLE_TIMEOUT"),
			TLSEnabled:   viper.GetBool("SERVER_TLS_ENABLED"),
			TLSCertFile:  viper.GetString("SERVER_TLS_CERT_FILE"),
			TLSKeyFile:   viper.GetString("SERVER_TLS_KEY_FILE"),
		},
		Database: DatabaseConfig{
			Host:         viper.GetString("DB_HOST"),
			Port:         viper.GetString("DB_PORT"),
			User:         viper.GetString("DB_USER"),
			Password:     viper.GetString("DB_PASSWORD"),
			DBName:       viper.GetString("DB_NAME"),
			SSLMode:      viper.GetString("DB_SSLMODE"),
			MaxOpenConns: viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns: viper.GetInt("DB_MAX_IDLE_CONNS"),
			MaxLifetime:  viper.GetDuration("DB_MAX_LIFETIME"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:          viper.GetString("JWT_SECRET"),
			AccessTokenTTL:  viper.GetDuration("JWT_ACCESS_TOKEN_TTL"),
			RefreshTokenTTL: viper.GetDuration("JWT_REFRESH_TOKEN_TTL"),
			Issuer:          viper.GetString("JWT_ISSUER"),
		},
		Cookie: CookieConfig{
			Domain:   viper.GetString("COOKIE_DOMAIN"),
			Secure:   viper.GetBool("COOKIE_SECURE"),
			SameSite: viper.GetString("COOKIE_SAMESITE"),
			Path:     viper.GetString("COOKIE_PATH"),
		},
		Log: LogConfig{
			Level:  viper.GetString("LOG_LEVEL"),
			Format: viper.GetString("LOG_FORMAT"),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseStringSlice("CORS_ALLOWED_ORIGINS"),
		},
		SMTP: SMTPConfig{
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetInt("SMTP_PORT"),
			Username: viper.GetString("SMTP_USERNAME"),
			Password: viper.GetString("SMTP_PASSWORD"),
			From:     viper.GetString("SMTP_FROM"),
		},
		OAuth: OAuthConfig{
			GoogleClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
			GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
			GoogleRedirectURL:  viper.GetString("GOOGLE_REDIRECT_URL"),
		},
		Telegram: TelegramConfig{
			BotToken: viper.GetString("BOT_TOKEN"),
			ChatID:   viper.GetString("TELEGRAM_CHAT_ID"),
			ProxyURL: viper.GetString("TELEGRAM_PROXY"),
		},
	}

	return cfg, nil
}

func validateSecuritySettings() error {
	// Validate JWT secret
	jwtSecret := viper.GetString("JWT_SECRET")
	if jwtSecret == "" {
		return fmt.Errorf("JWT_SECRET must be set and cannot be empty")
	}
	if len(jwtSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}
	if jwtSecret == "change-me-to-a-strong-secret-in-production" || jwtSecret == "change-me-in-production" {
		return fmt.Errorf("JWT_SECRET must be changed from default value")
	}

	// Validate database SSL mode in production
	if viper.GetString("APP_ENV") == "production" {
		sslMode := viper.GetString("DB_SSLMODE")
		if sslMode == "disable" {
			return fmt.Errorf("DB_SSLMODE cannot be 'disable' in production")
		}
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("APP_NAME", "luxbiss-server")
	viper.SetDefault("APP_ENV", "dev")

	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("SERVER_READ_TIMEOUT", "15s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "15s")
	viper.SetDefault("SERVER_IDLE_TIMEOUT", "60s")
	viper.SetDefault("SERVER_TLS_ENABLED", false)
	viper.SetDefault("SERVER_TLS_CERT_FILE", "")
	viper.SetDefault("SERVER_TLS_KEY_FILE", "")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "luxbiss")
	viper.SetDefault("DB_SSLMODE", "require")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 5)
	viper.SetDefault("DB_MAX_LIFETIME", "5m")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("JWT_SECRET", "")
	viper.SetDefault("JWT_ACCESS_TOKEN_TTL", "15m")
	viper.SetDefault("JWT_REFRESH_TOKEN_TTL", "168h")
	viper.SetDefault("JWT_ISSUER", "luxbiss-server")

	viper.SetDefault("COOKIE_DOMAIN", "")
	viper.SetDefault("COOKIE_SECURE", false)
	viper.SetDefault("COOKIE_SAMESITE", "Lax")
	viper.SetDefault("COOKIE_PATH", "/")

	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("LOG_FORMAT", "console")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", []string{})

	viper.SetDefault("SMTP_HOST", "smtp.gmail.com")
	viper.SetDefault("SMTP_PORT", 587)
	viper.SetDefault("SMTP_USERNAME", "")
	viper.SetDefault("SMTP_PASSWORD", "")
	viper.SetDefault("SMTP_FROM", "noreply@luxbiss.com")

	viper.SetDefault("GOOGLE_CLIENT_ID", "")
	viper.SetDefault("GOOGLE_CLIENT_SECRET", "")
	viper.SetDefault("GOOGLE_REDIRECT_URL", "")

	viper.SetDefault("BOT_TOKEN", "")
	viper.SetDefault("TELEGRAM_CHAT_ID", "")
	viper.SetDefault("TELEGRAM_PROXY", "")
}

func parseStringSlice(key string) []string {
	val := viper.GetString(key)
	if val == "" {
		return []string{}
	}

	// Viper.GetStringSlice doesn't always handle env vars correctly
	// if they are not explicitly typed. We manually split by comma.
	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
