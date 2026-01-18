package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"product-api/internal/database"
	"product-api/pkg/models"
	"strconv"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var item models.Product

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO products (name, description, price) 
	VALUES ($1, $2, $3) 
	RETURNING id
	`

	err := database.DB.QueryRow(query, item.Name, item.Description, item.Price).Scan(&item.ID)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
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

	var item models.Product

	query := `SELECT id, name, description, price
	FROM products
	WHERE id = $1
	`
	err = database.DB.QueryRow(query, id).Scan(&item.ID, &item.Name, &item.Description, &item.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error:"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
