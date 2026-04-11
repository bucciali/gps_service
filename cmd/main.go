package main

import (
	"context"
	"gps_service/internal/auth"
	"gps_service/internal/cache"
	"gps_service/internal/config"
	"gps_service/internal/db"
	"gps_service/internal/router"
	"log"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	pool, err := db.NewPostgresPool(context.Background(), cfg.DataBaseUrl)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	rdb := cache.NewRedisClient("localhost:6379")
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	pointsCache := cache.NewPointsCache(rdb)
	kiosksCache := cache.NewKioskCache(rdb)

	jm := auth.NewJWTManager(cfg.JwtSecret)
	r := router.NewRouter(jm, pool, pointsCache, kiosksCache)

	srv := &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}

	log.Println("server stopped")
}
