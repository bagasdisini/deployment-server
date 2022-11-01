package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	ShowUsers() ([]models.User, error)
	ShowAdmins() ([]models.Admin, error)
	GetUserByIDUser(ID int) (models.User, error)
	GetUserByIDAdmin(ID int) (models.Admin, error)
	CreateUserUser(user models.User) (models.User, error)
	CreateUserAdmin(admin models.Admin) (models.Admin, error)
	UpdateUser(user models.User, ID int) (models.User, error)
	UpdateAdmin(admin models.Admin, ID int) (models.Admin, error)
	DeleteUser(user models.User, ID int) (models.User, error)
	DeleteAdmin(admin models.Admin, ID int) (models.Admin, error)
}

func RepositoryUser(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ShowUsers() ([]models.User, error) {
	var user []models.User
	err := r.db.Find(&user).Error

	return user, err
}

func (r *repository) ShowAdmins() ([]models.Admin, error) {
	var Admin []models.Admin
	err := r.db.Find(&Admin).Error

	return Admin, err
}

func (r *repository) GetUsers(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error

	return user, err
}

func (r *repository) GetAdmins(ID int) (models.Admin, error) {
	var admin models.Admin
	err := r.db.First(&admin, ID).Error

	return admin, err
}

func (r *repository) CreateUserUser(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error

	return user, err
}

func (r *repository) CreateUserAdmin(admin models.Admin) (models.Admin, error) {
	err := r.db.Create(&admin).Error

	return admin, err
}

func (r *repository) UpdateUser(user models.User, ID int) (models.User, error) {
	err := r.db.Model(&user).Where("id=?", ID).Updates(&user).Error

	return user, err
}

func (r *repository) UpdateAdmin(admin models.Admin, ID int) (models.Admin, error) {
	err := r.db.Model(&admin).Where("id=?", ID).Updates(&admin).Error

	return admin, err
}

func (r *repository) DeleteUser(user models.User, ID int) (models.User, error) {
	err := r.db.Delete(&user).Error

	return user, err
}

func (r *repository) DeleteAdmin(admin models.Admin, ID int) (models.Admin, error) {
	err := r.db.Delete(&admin).Error

	return admin, err
}
