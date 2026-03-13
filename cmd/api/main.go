package main

import (
	"context"
	"goCachedAPI/internal/cache"
	"goCachedAPI/internal/config"
	"goCachedAPI/internal/db"
	"goCachedAPI/internal/handlers"
	"goCachedAPI/internal/repository"
	"goCachedAPI/internal/routes"
	"goCachedAPI/internal/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>
func main() {
	cfg := config.Load()

	database := db.New(cfg.SQLiteDSN)

	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	productRepo := repository.NewProductRepository(database)
	productCache := cache.NewProductCache(
		redisClient,
		cfg.ProductTTL,
		cfg.RecentListKey,
		cfg.RecentLimit,
	)

	productService := service.NewProductService(productRepo, productCache)
	productHandler := handlers.NewProductHandler(productService)

	router := gin.Default()
	routes.Register(router, productHandler)

	log.Printf("server running on %s", cfg.AppPort)
	if err := router.Run(cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
