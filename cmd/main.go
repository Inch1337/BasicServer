package main

import (
	"fmt"
	"log"
	"net/http"
	"product-api/internal/config"
	"product-api/internal/database"
	"product-api/internal/handlers"

	_ "github.com/lib/pq"
)

func main() {
	config.Load()
	cfg := config.New()

	database.InitDb(cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/products", handlers.ProductsHandler)
	mux.HandleFunc("/products/", handlers.SingleProductHandler)

	fmt.Printf("Start server on %s\n", cfg.ServerPort)

	if err := http.ListenAndServe(cfg.ServerPort, mux); err != nil {
		log.Fatal(err)
	}
}
