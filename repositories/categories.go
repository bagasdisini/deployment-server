package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	ShowCategory() ([]models.Category, error)
}

func RepositoryCategory(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ShowCategory() ([]models.Category, error) {
	var Category []models.Category
	err := r.db.Find(&Category).Error

	return Category, err
}
