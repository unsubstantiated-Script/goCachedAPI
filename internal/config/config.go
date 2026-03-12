package config

import "time"

type Config struct {
	AppPort       string
	RedisAddr     string
	SQLiteDSN     string
	ProductTTL    time.Duration
	RecentListKey string
	RecentLimit   int64
}

func Load() Config {
	return Config{
		AppPort:       ":8080",
		RedisAddr:     "localhost:6379",
		SQLiteDSN:     "app.db",
		ProductTTL:    time.Minute,
		RecentListKey: "recent_products",
		RecentLimit:   10,
	}
}
