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
	UpdateTransaction(status string, ID string) error
	GetOneTransaction(ID string) (models.TransactionUser, error)
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

func (r *repository) GetOneTransaction(ID string) (models.TransactionUser, error) {
	var transaction models.TransactionUser
	err := r.db.Preload("Product").Preload("Product.User").Preload("Buyer").Preload("Seller").First(&transaction, "id = ?", ID).Error

	return transaction, err
}

func (r *repository) UpdateTransaction(status string, ID string) error {
	var transaction models.TransactionUser
	r.db.Preload("Product").First(&transaction, ID)

	if status != transaction.Status && status == "success" {
		var product models.Product
		r.db.First(&product, transaction.Product)
		r.db.Save(&product)
	}

	transaction.Status = status

	err := r.db.Save(&transaction).Error

	return err
}
