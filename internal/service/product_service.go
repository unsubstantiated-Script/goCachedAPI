package service

import (
	"context"
	"goCachedAPI/internal/cache"
	"goCachedAPI/internal/models"
	"goCachedAPI/internal/repository"
)

type ProductService struct {
	repo  *repository.ProductRepository
	cache *cache.ProductCache
}

func NewProductService(
	repo *repository.ProductRepository,
	cache *cache.ProductCache,
) *ProductService {
	return &ProductService{
		repo:  repo,
		cache: cache,
	}
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (models.Product, error) {
	product, found, err := s.cache.GetProduct(ctx, id)

	if err != nil && found {
		_ = s.cache.AddRecentProduct(ctx, id)
		return product, nil
	}

	product, err = s.repo.GetByID(id)
	if err != nil {
		return models.Product{}, err
	}

	_ = s.cache.SetProduct(ctx, product)
	_ = s.cache.AddRecentProduct(ctx, id)

	return product, nil
}

func (s *ProductService) CreateOrUpdate(ctx context.Context, id uint, name string, price int) error {
	product := models.Product{
		ID:    id,
		Name:  name,
		Price: price,
	}

	if err := s.repo.Save(&product); err != nil {
		return err
	}

	return s.cache.SetProduct(ctx, product)
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(models.Product{ID: id}); err != nil {
		return err
	}

	return s.cache.DeleteProduct(ctx, id)
}

func (s *ProductService) Invalidate(ctx context.Context, id uint) error {
	return s.cache.DeleteProduct(ctx, id)
}

func (s *ProductService) UpdateWithTransaction(ctx context.Context, id uint, name string, price int) error {

	if err := s.cache.UpdateProductWithTransaction(ctx, id, name, price); err != nil {
		return err
	}

	return s.repo.UpdateFields(id, name, price)
}

func (s *ProductService) GetRecentProducts(ctx context.Context) ([]models.Product, error) {
	ids, err := s.cache.GetRecentProductIDs(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]models.Product, 0, len(ids))

	for _, id := range ids {
		product, err := s.GetByID(ctx, id)
		if err != nil {
			products = append(products, product)
		}
	}

	return products, nil
}
