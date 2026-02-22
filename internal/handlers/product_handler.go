package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"product-test/internal/apierr"
	"product-test/internal/models"
	"product-test/internal/service"
	"strconv"
)

type ProductHandler struct {
	service service.ProductService
	log     *slog.Logger
}

func NewProductHandler(svc service.ProductService, log *slog.Logger) *ProductHandler {
	if log == nil {
		log = slog.Default()
	}
	return &ProductHandler{service: svc, log: log}
}

func (h *ProductHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /products", h.getAll)
	mux.HandleFunc("POST /products", h.create)
	mux.HandleFunc("GET /products/{id}", h.getByID)
	mux.HandleFunc("PUT /products/{id}", h.update)
	mux.HandleFunc("DELETE /products/{id}", h.delete)
}

func (h *ProductHandler) getAll(w http.ResponseWriter, r *http.Request) {
	limit, offset := parseLimitOffset(r)
	products, err := h.service.GetAllProducts(r.Context(), limit, offset)
	if err != nil {
		h.log.Error("get all products", "error", err)
		apierr.Internal(w)
		return
	}
	h.writeJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		apierr.BadRequest(w, "invalid JSON")
		return
	}
	p.ID = 0
	if err := h.service.CreateProduct(r.Context(), &p); err != nil {
		if errors.Is(err, service.ErrValidation) {
			apierr.BadRequest(w, err.Error())
			return
		}
		h.log.Error("create product", "error", err)
		apierr.Internal(w)
		return
	}
	h.writeJSON(w, http.StatusCreated, p)
}

func (h *ProductHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	p, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			apierr.NotFound(w, "product not found")
			return
		}
		h.log.Error("get product by id", "id", id, "error", err)
		apierr.Internal(w)
		return
	}
	h.writeJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		apierr.BadRequest(w, "invalid JSON")
		return
	}
	p.ID = id
	if err := h.service.UpdateProduct(r.Context(), &p); err != nil {
		if errors.Is(err, service.ErrValidation) {
			apierr.BadRequest(w, err.Error())
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			apierr.NotFound(w, "product not found")
			return
		}
		h.log.Error("update product", "id", id, "error", err)
		apierr.Internal(w)
		return
	}
	h.writeJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	if err := h.service.DeleteProduct(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			apierr.NotFound(w, "product not found")
			return
		}
		h.log.Error("delete product", "id", id, "error", err)
		apierr.Internal(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseLimitOffset(r *http.Request) (limit, offset int) {
	limit = 100
	offset = 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
			if limit > 500 {
				limit = 500
			}
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return limit, offset
}

func parseID(w http.ResponseWriter, r *http.Request, param string) (int, bool) {
	idStr := r.PathValue(param)
	if idStr == "" {
		apierr.BadRequest(w, "invalid product id")
		return 0, false
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		apierr.BadRequest(w, "invalid product id")
		return 0, false
	}
	return id, true
}

func (h *ProductHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		h.log.Error("encode response", "error", err)
	}
}
