package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mymail/worker/internal/config"
	"github.com/mymail/worker/internal/processor"
	"github.com/mymail/worker/internal/storage"
)

func main() {
	var cfgPath string
	flag.StringVar(&cfgPath, "config", "", "Path to config file")
	flag.Parse()

	cfg := config.Load(cfgPath)

	// Initialize storage
	db, err := storage.NewPostgres(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redis, err := storage.NewRedis(cfg.Redis.URL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Create processor
	proc := processor.New(db, redis, cfg)

	// Start workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start processing loop
	go proc.Start(ctx)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	cancel()

	// Wait for workers to finish
	time.Sleep(2 * time.Second)
}
