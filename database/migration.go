package database

import (
	"backend/models"
	"backend/pkg/mysql"
	"fmt"
)

func RunMigration() {
	err := mysql.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.TransactionUser{}, &models.Admin{}, &models.Category{})

	if err != nil {
		fmt.Println(err)
		panic("Migration Failed!")
	}

	fmt.Println("Migration Successful!")
}
