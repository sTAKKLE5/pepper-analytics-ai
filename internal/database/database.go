package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"pepper-analytics-ai/internal/utils"
	"time"
)

type Config struct {
	Host                    string
	Port                    int
	User                    string
	Password                string
	DBName                  string
	SSLMode                 string
	MaxOpenConns            int
	MaxIdleConns            int
	ConnMaxLifetimeMinutes  int
	ConnectionMaxRetries    int
	ConnectionRetryDelaySec int
}

func NewDefaultConfig() *Config {
	return &Config{
		Host:                    utils.GetRequiredEnv("DB_HOST"),
		Port:                    utils.GetRequiredEnvAsInt("DB_PORT"),
		User:                    utils.GetRequiredEnv("DB_USER"),
		Password:                utils.GetRequiredEnv("DB_PASSWORD"),
		DBName:                  utils.GetRequiredEnv("DB_NAME"),
		SSLMode:                 utils.GetRequiredEnv("DB_SSLMODE"),
		MaxOpenConns:            utils.GetRequiredEnvAsInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns:            utils.GetRequiredEnvAsInt("DB_MAX_IDLE_CONNS"),
		ConnMaxLifetimeMinutes:  utils.GetRequiredEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES"),
		ConnectionMaxRetries:    utils.GetRequiredEnvAsInt("DB_MAX_RETRIES"),
		ConnectionRetryDelaySec: utils.GetRequiredEnvAsInt("DB_RETRY_DELAY_SECONDS"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func NewConnection(config *Config) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	// Try to connect with retries
	for i := 0; i < config.ConnectionMaxRetries; i++ {
		db, err = sqlx.Connect("postgres", config.DSN())
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v",
			i+1, config.ConnectionMaxRetries, err)

		if i < config.ConnectionMaxRetries-1 {
			time.Sleep(time.Duration(config.ConnectionRetryDelaySec) * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w",
			config.ConnectionMaxRetries, err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetimeMinutes) * time.Minute)

	// Verify connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error verifying database connection: %w", err)
	}

	// Log configuration
	log.Printf("Database pool configured with: max_open=%d, max_idle=%d, lifetime=%dm",
		config.MaxOpenConns, config.MaxIdleConns, config.ConnMaxLifetimeMinutes)

	return db, nil
}
