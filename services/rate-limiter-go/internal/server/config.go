package server

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"ratelimiter"
)

// Config is everything the server reads from the environment.
type Config struct {
	Port      string
	RateLimit ratelimiter.Config
}

// ConfigFromEnv reads Config from environment variables, using defaults for
// any that are unset.
//
//	PORT                                8080
//	RATE_LIMIT_CAPACITY                 100
//	RATE_LIMIT_INTERVAL_SECONDS         60
//	RATE_LIMIT_REFILL_AMOUNT            10
//	RATE_LIMIT_REFILL_INTERVAL_SECONDS  6
func ConfigFromEnv() (Config, error) {
	capacity, err := envInt("RATE_LIMIT_CAPACITY", 100)
	if err != nil {
		return Config{}, err
	}
	interval, err := envInt("RATE_LIMIT_INTERVAL_SECONDS", 60)
	if err != nil {
		return Config{}, err
	}
	refillAmount, err := envInt("RATE_LIMIT_REFILL_AMOUNT", 10)
	if err != nil {
		return Config{}, err
	}
	refillInterval, err := envInt("RATE_LIMIT_REFILL_INTERVAL_SECONDS", 6)
	if err != nil {
		return Config{}, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Port: port,
		RateLimit: ratelimiter.Config{
			Capacity:       capacity,
			Interval:       time.Duration(interval) * time.Second,
			RefillAmount:   refillAmount,
			RefillInterval: time.Duration(refillInterval) * time.Second,
		},
	}, nil
}

func envInt(name string, fallback int) (int, error) {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s=%q is not an integer", name, raw)
	}
	return v, nil
}
