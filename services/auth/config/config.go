package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config agrupa toda la configuración del servicio.
// Se carga desde variables de entorno — nunca hardcodeada.
type Config struct {
	App    AppConfig
	HTTP   HTTPConfig
	DB     DBConfig
	JWT    JWTConfig
	Google GoogleConfig
	Email  EmailConfig
	Legal  LegalConfig
}

type AppConfig struct {
	Name     string
	Env      string // development | production
	LogLevel string // debug | info | warn | error
}

type HTTPConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// DSN retorna el connection string para pgx.
func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

type JWTConfig struct {
	AccessSecret       string
	RefreshSecret      string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromAddress  string
	FromName     string
}

type LegalConfig struct {
	ConsentPrivacyVersion      string
	ConsentDataCampaignVersion string
}

// Load carga la configuración desde variables de entorno.
// En desarrollo lee el .env automáticamente si existe.
// En producción las variables vienen del entorno del sistema — no hay .env.
// Falla rápido (panic) si alguna variable crítica no está definida.
func Load() Config {
	// Carga el .env si existe — en producción simplemente no hace nada
	_ = godotenv.Load()

	return Config{
		App: AppConfig{
			Name:     getEnv("APP_NAME", "auth-service"),
			Env:      getEnv("APP_ENV", "development"),
			LogLevel: getEnv("LOG_LEVEL", "debug"),
		},
		HTTP: HTTPConfig{
			Port:         getEnvInt("HTTP_PORT", 8080),
			ReadTimeout:  time.Duration(getEnvInt("HTTP_READ_TIMEOUT_SEC", 10)) * time.Second,
			WriteTimeout: time.Duration(getEnvInt("HTTP_WRITE_TIMEOUT_SEC", 10)) * time.Second,
		},
		DB: DBConfig{
			Host:     mustGetEnv("DB_HOST"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     mustGetEnv("DB_USER"),
			Password: mustGetEnv("DB_PASSWORD"),
			Name:     mustGetEnv("DB_NAME"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			AccessSecret:       mustGetEnv("JWT_ACCESS_SECRET"),
			RefreshSecret:      mustGetEnv("JWT_REFRESH_SECRET"),
			AccessTokenExpiry:  time.Duration(getEnvInt("JWT_ACCESS_EXPIRY_MIN", 15)) * time.Minute,
			RefreshTokenExpiry: time.Duration(getEnvInt("JWT_REFRESH_EXPIRY_DAYS", 7)) * 24 * time.Hour,
		},
		Google: GoogleConfig{
			ClientID:     mustGetEnv("GOOGLE_CLIENT_ID"),
			ClientSecret: mustGetEnv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  mustGetEnv("GOOGLE_REDIRECT_URL"),
		},
		Email: EmailConfig{
			SMTPHost:     mustGetEnv("SMTP_HOST"),
			SMTPPort:     getEnvInt("SMTP_PORT", 587),
			SMTPUser:     mustGetEnv("SMTP_USER"),
			SMTPPassword: mustGetEnv("SMTP_PASSWORD"),
			FromAddress:  getEnv("EMAIL_FROM_ADDRESS", "no-reply@pc-store.com"),
			FromName:     getEnv("EMAIL_FROM_NAME", "PC Store"),
		},
		Legal: LegalConfig{
			ConsentPrivacyVersion:      mustGetEnv("CONSENT_PRIVACY_VERSION"),
			ConsentDataCampaignVersion: mustGetEnv("CONSENT_DATA_CAMPAIGN_VERSION"),
		},
	}
}

// -------------------------
// Helpers
// -------------------------

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// mustGetEnv falla si la variable no está definida — para configuración crítica.
func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("config: environment variable %q is required but not set", key))
	}
	return v
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
