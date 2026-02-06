package service

import (
	"context"
	"errors"
	"product-test/internal/models"
	"product-test/internal/repository"
)

type ProductService interface {
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id int) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *productService) CreateProduct(ctx context.Context, product *models.Product) error {
	if product.Price < 0 {
		return errors.New("price cannot be negative")
	}
	return s.repo.Create(ctx, product)
}

func (s *productService) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) UpdateProduct(ctx context.Context, product *models.Product) error {
	if product.Price < 0 {
		return errors.New("price cannot be negative")
	}
	return s.repo.Update(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
