package models

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	FirstName string `gorm:"size:100;not null"`
	LastName  string `gorm:"size:100;not null"`
	Email     string `gorm:"size:255;not null;unique"`
	Phone     string `gorm:"size:20;not null"`
	Address   string `gorm:"size:255"`
	OAuthID   string `gorm:"column:oauth_id;size:255;unique"`
}

// 	gorm.Model This is an embedded struct provided by GORM. It includes the following fields automatically:
//ID        uint      // primary key
//CreatedAt time.Time
//UpdatedAt time.Time
//DeletedAt gorm.DeletedAt `gorm:"index"` // for soft deletes
