package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/mymail/smtp/src/config"
	"github.com/mymail/smtp/src/handler"
	"github.com/mymail/smtp/src/ratelimit"
	"github.com/mymail/smtp/src/storage"
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

	minio, err := storage.NewMinIO(cfg.MinIO)
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	// Initialize rate limiter
	rateLimiter := ratelimit.New(redis)

	// Create SMTP backend
	backend := handler.NewBackend(db, redis, minio, rateLimiter, cfg)

	// Create SMTP server
	s := smtp.NewServer(backend)
	s.Addr = fmt.Sprintf("%s:%d", cfg.SMTP.Host, cfg.SMTP.Port)
	s.Domain = cfg.SMTP.Domain
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = cfg.SMTP.MaxMessageSize
	s.MaxRecipients = 50
	s.AllowInsecureAuth = false

	// TLS configuration
	if cfg.TLS.Enabled {
		cert, err := tls.LoadX509KeyPair(cfg.TLS.CertFile, cfg.TLS.KeyFile)
		if err != nil {
			log.Fatalf("Failed to load TLS certificate: %v", err)
		}
		s.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			ServerName:   cfg.SMTP.Domain,
		}
	}

	// Start server
	log.Printf("Starting SMTP server on %s", s.Addr)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("SMTP server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down SMTP server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}
}
