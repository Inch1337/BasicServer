package service

import (
	"context"
	"errors"
	"fmt"
	"product-test/internal/models"
	"product-test/internal/repository"
)

const (
	MaxNameLength        = 500
	MaxDescriptionLength = 2000
)

type ProductService interface {
	GetAllProducts(ctx context.Context, limit, offset int) ([]models.Product, error)
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

func (s *productService) GetAllProducts(ctx context.Context, limit, offset int) ([]models.Product, error) {
	return s.repo.GetAll(ctx, limit, offset)
}

func validateProduct(p *models.Product) error {
	if p.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	if len(p.Name) > MaxNameLength {
		return fmt.Errorf("%w: name must be at most %d characters", ErrValidation, MaxNameLength)
	}
	if len(p.Description) > MaxDescriptionLength {
		return fmt.Errorf("%w: description must be at most %d characters", ErrValidation, MaxDescriptionLength)
	}
	if p.Price < 0 {
		return fmt.Errorf("%w: price cannot be negative", ErrValidation)
	}
	return nil
}

func (s *productService) CreateProduct(ctx context.Context, product *models.Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}
	return s.repo.Create(ctx, product)
}

func (s *productService) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return p, err
}

func (s *productService) UpdateProduct(ctx context.Context, product *models.Product) error {
	if err := validateProduct(product); err != nil {
		return err
	}
	return s.repo.Update(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}
	return err
}
