// Package main Product API server.
//
//	@title			Product API
//	@version		1.0
//	@description	REST API for product management
//	@host			localhost:8081
//	@BasePath		/
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"product-test/internal/config"
	"product-test/internal/database"
	"product-test/internal/handlers"
	"product-test/internal/repository"
	"product-test/internal/service"

	_ "product-test/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	config.LoadEnv()
	cfg, err := config.New()
	if err != nil {
		slog.Error("invalid config", "error", err)
		os.Exit(1)
	}

	db, err := database.InitDB(cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)
	if err != nil {
		slog.Error("database init", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	productHandler := handlers.NewProductHandler(productService, logger)

	mux := http.NewServeMux()
	productHandler.RegisterRoutes(mux)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: mux,
	}

	go func() {
		logger.Info("server started", "addr", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server shutdown", "error", err)
	}
	logger.Info("server stopped")
}
