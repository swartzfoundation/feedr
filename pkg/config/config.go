package config

import (
	"context"
	"os"

	"log/slog"

	"github.com/sethvargo/go-envconfig"
)

var (
	Config *config
)

type SessionConfig struct {
	// SessionKey is the key used to sign the session cookie
	SessionCookieName   string `env:"SESSION_COOKIE_NAME"`
	SessionCookieDomain string `env:"SESSION_COOKIE_DOMAIN,required"`
}

// DBConfig contains the configuration for the database.
type DBConfig struct {
	// PostgresHost is the host of the Postgres server
	PostgresHost string `env:"POSTGRES_HOST,default=localhost"`
	// PostgresPort is the port of the Postgres server
	PostgresPort string `env:"POSTGRES_PORT,default=5432"`
	// PostgresDatabase is the name of the Postgres database
	DatabaseName string `env:"POSTGRES_DATABASE,default=feedr"`
	// PostgresUsername is the username of the Postgres server
	PostgresUsername string `env:"POSTGRES_USERNAME,required"`
	// PostgresPassword is the password of the Postgres server
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	// PostgresSSLMode is the sslmode of the Postgres server
	PostgresSSLMode string `env:"POSTGRES_SSLMODE,default=prefer"`
	DSN             string `env:"DSN"`
}

// EmailConfig contains the configuration for the email server.
type EmailConfig struct {
	// SMTPHost is the host of the SMTP server
	SMTPHost string `env:"SMTP_HOST"`
	// SMTPPort is the port of the SMTP server
	SMTPPort string `env:"SMTP_PORT,default=587"`
	// SMTPAuth is whether the SMTP server requires authentication
	SMTPAuth bool `env:"SMTP_TLS,default=false"`
	// SMTPUsername is the username of the SMTP server
	SMTPUsername string `env:"SMTP_USERNAME"`
	// SMTPPassword is the password of the SMTP server
	SMTPPassword string `env:"SMTP_PASSWORD"`
	// FromAddress is the email address that emails will be sent from
	FromAddress string `env:"FROM_ADDRESS"`
	FromName    string `env:"FROM_NAME"`
}

type OpenAIConfig struct {
	// APIKey is the API key for the OpenAI API
	APIKey string `env:"OPENAI_API_KEY"`
}

type config struct {
	DEBUG           bool     `env:"DEBUG,default=false"`
	PORT            string   `env:"PORT,default=8000"`
	ALLOWED_ORIGINS []string `env:"ALLOWED_ORIGINS,default=*"`
	Session         SessionConfig

	Email  EmailConfig
	OpenAI OpenAIConfig
	DB     DBConfig
}

func Load() *config {
	c := config{}
	if err := envconfig.Process(context.TODO(), &c); err != nil {
		slog.Error("failed to load env", "error", err)
		os.Exit(1)
	}
	Config = &c
	return &c
}
