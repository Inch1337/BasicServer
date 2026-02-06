package main

import (
	"log"
	"net/http"
	"product-test/internal/config"
	"product-test/internal/database"
	"product-test/internal/handlers"
	"product-test/internal/repository"
	"product-test/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	config.Load()
	cfg := config.New()

	database.InitDb(cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)

	productRepo := repository.NewProductRepository(database.DB)

	productService := service.NewProductService(productRepo)

	productHandler := handlers.NewProductHandler(productService)

	mux := http.NewServeMux()
	productHandler.RegisterRoutes(mux)

	log.Printf("Start server on %s\n", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, mux); err != nil {
		log.Fatal(err)
	}
}
