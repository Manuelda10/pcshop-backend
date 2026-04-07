package config

import "os"

type Config struct {
	App   AppConfig
	Email EmailConfig
}

type AppConfig struct {
	Name     string
	Env      string
	LogLevel string
}

type EmailConfig struct {
	// FromAddress es el email verificado en SES desde el cual se envían los emails.
	// Debe estar verificado en SES (dominio o dirección individual).
	FromAddress string
}

func Load() Config {
	return Config{
		App: AppConfig{
			Name:     getEnv("APP_NAME", "PCShop"),
			Env:      getEnv("APP_ENV", "development"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		Email: EmailConfig{
			FromAddress: getEnv("EMAIL_FROM_ADDRESS", "noreply@pcshop.com"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
