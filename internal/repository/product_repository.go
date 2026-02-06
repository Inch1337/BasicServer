package repository

import (
	"context"
	"database/sql"
	"errors"
	"product-test/internal/models"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `SELECT id, name, description, price FROM products`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepo) Create(ctx context.Context, p *models.Product) error {
	query := `INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, p.Name, p.Description, p.Price).Scan(&p.ID)
}

func (r *productRepo) GetByID(ctx context.Context, id int) (*models.Product, error) {
	query := `SELECT id, name, description, price FROM products WHERE id = $1`
	var p models.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) Update(ctx context.Context, p *models.Product) error {
	query := `UPDATE products SET name=$1, description=$2, price=$3 WHERE id=$4`
	res, err := r.db.ExecContext(ctx, query, p.Name, p.Description, p.Price, p.ID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *productRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}
	return nil
}
