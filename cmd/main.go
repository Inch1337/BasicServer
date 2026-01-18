package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"product-api/internal/config"
	"product-api/internal/database"
	"strconv"

	_ "github.com/lib/pq"
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

var (
	products  = []Product{}
	idCounter = 1
)

var DB *sql.DB

func main() {
	config.Load()
	cfg := config.New()

	database.InitDb(cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/products", ProductsHandler)
	mux.HandleFunc("/products/", GetProductHandler)

	fmt.Printf("Start server on %s\n", cfg.ServerPort)

	if err := http.ListenAndServe(cfg.ServerPort, mux); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "Welcome to Product API")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "Service running")
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)

	case "POST":
		var p Product

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "Bad JSON", http.StatusBadRequest)
			return
		}

		p.ID = idCounter
		idCounter++

		products = append(products, p)

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(p)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/products/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, item := range products {
		if item.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)

			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}
