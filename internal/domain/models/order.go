package models

import (
	"errors"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

var ErrOrderNotFound = errors.New("order not found")

type Order struct {
	gorm.Model
	CustomerID uint        `gorm:"not null"`
	Customer   Customer    `gorm:"foreignkey:CustomerID"`
	Status     OrderStatus `gorm:"type:varchar(20);default:'pending'"`
	Total      float64     `gorm:"type:decimal(10,2);not null"`
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Product   Product `gorm:"foreignkey:ProductID"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"type:decimal(10,2);not null"`
}
