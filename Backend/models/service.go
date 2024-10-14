package models

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	ProviderID   uint    `gorm:"not null"`
	Name         string  `json:"name" gorm:"not null"`
	Description  string  `json:"description" gorm:"not null"`
	Price        float64 `json:"price" gorm:"not null"`        
	Rating       float64 `json:"rating"`       // Average rating of this service
	Availability string  `json:"availability" gorm:"not null" ` // "Mon-Fri 9 AM - 5 PM"
	Duration     int     `json:"duration"`     // Duration of the service in minutes
	Visibility   bool    `json:"visible"`
}

type Review struct {
	gorm.Model
	ServiceID    uint      `gorm:"not null"` // Foreign key to the Service
	CustomerID   uint      `gorm:"not null"` // Foreign key to the Customer
	Rating       float64   `json:"rating"`   // Rating given by the customer
	Comment      string    `json:"comment"`  // Review comment
	ReviewDate   time.Time `json:"review_date"`
	CustomerName string    `json:"customer_name"` // For displaying the name of the customer who reviewed
}
