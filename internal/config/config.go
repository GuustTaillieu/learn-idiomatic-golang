package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port         string
	DatabaseURL  string
	WorkerCount  int
	RateLimit    float64
	RateBurst    int
	ShutdownTime time.Duration
}

func Load(args ...string) (*Config, error) {
	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "file:orders.db?cache=shared&mode=rwc"),
		WorkerCount:  getEnvInt("WORKER_COUNT", 3),
		RateLimit:    getEnvFloat("RATE_LIMIT", 5.0),
		RateBurst:    getEnvInt("RATE_BURST", 10),
		ShutdownTime: 5 * time.Second,
	}

	// Use an isolated FlagSet per Load call so tests don't panic with 'flag redefined'
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.StringVar(&cfg.Port, "port", cfg.Port, "Port to run the server on")
	fs.StringVar(&cfg.DatabaseURL, "db", cfg.DatabaseURL, "Database connection URL")
	fs.IntVar(&cfg.WorkerCount, "workers", cfg.WorkerCount, "Number of workers to process tasks")
	fs.Float64Var(&cfg.RateLimit, "rate-limit", cfg.RateLimit, "Rate limit for requests per second")
	fs.IntVar(&cfg.RateBurst, "rate-burst", cfg.RateBurst, "Burst size for rate limiting")
	fs.DurationVar(&cfg.ShutdownTime, "shutdown-time", cfg.ShutdownTime, "Graceful shutdown time")

	if len(args) > 0 {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.WorkerCount <= 0 {
		return errors.New("WORKER_COUNT must be > 0")
	}
	if c.Port == "" {
		return errors.New("PORT cannot be empty")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		var intValue int
		_, err := fmt.Sscanf(value, "%d", &intValue)
		if err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		var floatValue float64
		_, err := fmt.Sscanf(value, "%f", &floatValue)
		if err == nil {
			return floatValue
		}
	}
	return defaultValue
}
