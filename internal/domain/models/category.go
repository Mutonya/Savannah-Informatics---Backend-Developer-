package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string `gorm:"size:100;not null"`
	ParentID *uint
	Parent   *Category  `gorm:"foreignkey:ParentID"`
	Children []Category `gorm:"foreignkey:ParentID"`
	Products []Product  `gorm:"foreignkey:CategoryID"`
}
