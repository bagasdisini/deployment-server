package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	ShowTransactionUser() ([]models.TransactionUser, error)
	GetTransactionByIDUser(ID int) (models.TransactionUser, error)
	CreateTransactionUser(transaction models.TransactionUser) (models.TransactionUser, error)
	UpdateTransactionUser(transaction models.TransactionUser, ID int) (models.TransactionUser, error)
	DeleteTransactionUser(transaction models.TransactionUser, ID int) (models.TransactionUser, error)
}

func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ShowTransactionUser() ([]models.TransactionUser, error) {
	var transactions []models.TransactionUser
	err := r.db.Find(&transactions).Error

	return transactions, err
}

func (r *repository) GetTransactionByIDUser(ID int) (models.TransactionUser, error) {
	var transactions models.TransactionUser
	err := r.db.First(&transactions, ID).Error

	return transactions, err
}

func (r *repository) CreateTransactionUser(transaction models.TransactionUser) (models.TransactionUser, error) {
	err := r.db.Create(&transaction).Error

	return transaction, err
}

func (r *repository) UpdateTransactionUser(transaction models.TransactionUser, ID int) (models.TransactionUser, error) {
	err := r.db.Model(&transaction).Where("id=?", ID).Updates(&transaction).Error

	return transaction, err
}

func (r *repository) DeleteTransactionUser(transaction models.TransactionUser, ID int) (models.TransactionUser, error) {
	err := r.db.Delete(&transaction).Error

	return transaction, err
}
