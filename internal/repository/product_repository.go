package repository

import (
	"goCachedAPI/internal/models"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(id uint) (models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	return product, err
}

func (r *ProductRepository) Save(product *models.Product) error {
	return r.db.Save(&product).Error
}

func (r *ProductRepository) Delete(product models.Product) error {
	return r.db.Delete(&product).Error
}

func (r *ProductRepository) UpdateFields(id uint, name string, price int) error {
	return r.db.Model(&models.Product{}).
		Where("id = ?", id).
		Updates(models.Product{Name: name, Price: price}).
		Error
}
