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

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetAllProducts(w, r)
	case http.MethodPost:
		CreateProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func SingleProductHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProductByID(w, r)
	case http.MethodPut:
		updateProduct(w, r)
	case http.MethodDelete:
		deleteProduct(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, description, price
	FROM products
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		http.Error(w, "Database error:"+err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	allProducts := []models.Product{}

	for rows.Next() {
		var p models.Product

		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err != nil {
			http.Error(w, "Database scan error", http.StatusInternalServerError)
			return
		}

		allProducts = append(allProducts, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Database iteration error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allProducts)
}

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

func getProductByID(w http.ResponseWriter, r *http.Request) {
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

func updateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	query := `UPDATE products SET name=$1, description=$2, price=$3 WHERE id=$4`

	res, err := database.DB.Exec(query, p.Name, p.Description, p.Price, id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	p.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM products WHERE id = $1`

	res, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
