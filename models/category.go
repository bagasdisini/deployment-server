package models

type Category struct {
	ID   int    `json:"id" gorm:"primary_key:auto_increment;"`
	Name string `json:"name" form:"fullName" gorm:"type: varchar(255)"`
}
