package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	RedisAddr     string
	SQLiteDSN     string
	ProductTTL    time.Duration
	RecentListKey string
	RecentLimit   int64
}

func Load() Config {

	_ = godotenv.Load()

	ttlSeconds := mustInt("PRODUCT_TTL_SECONDS", 60)
	recentLimit := mustInt("RECENT_LIMIT", 10)

	cfg := Config{
		AppPort:       getEnv("APP_PORT", ":8080"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		SQLiteDSN:     getEnv("SQLITE_DSN", "app.db"),
		ProductTTL:    time.Duration(ttlSeconds) * time.Second,
		RecentListKey: getEnv("RECENT_LIST_KEY", "recent_products"),
		RecentLimit:   int64(recentLimit),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func mustInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("invalid int for %s, using fallback %d", key, fallback)
		return fallback
	}
	return i
}
