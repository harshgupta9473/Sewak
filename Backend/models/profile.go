package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomerProfile struct {
	gorm.Model
	UserID                  uint    `gorm:"not null"`
	FirstName               string  `json:"first_name" gorm:"not null"`
	LastName                string  `json:"last_name"`
	Bio                     string  `json:"bio"`
	ProfilePicture          string  `json:"profile_picture"`          // URL to the profile picture
	Latitude                float64 `json:"latitude" gorm:"not null"` // Geographic coordinates for location
	Longitude               float64 `json:"longitude" gorm:"not null"`
	City                    string  `json:"city"`                     // City or region
	PreferredRadius         float64 `json:"preferred_radius"`         // Distance within which they prefer services (in km)
	ServiceCategories       string  `json:"service_categories"`       // JSON or comma-separated string of preferred service categories
	NotificationPreferences string  `json:"notification_preferences"` // Example: JSON for email, SMS preferences
	ContactEmail            string  `json:"contact_email"`
	FavoriteProviders       string  `json:"favorite_providers"` // Example: JSON list of provider IDs           // JSON or separate table for reviews given by the customer

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ServiceProviderProfile struct {
	gorm.Model
	UserID         uint      `gorm:"not null"`
	Name           string    `json:"name" gorm:"not null"`
	Bio            string    `json:"bio" gorm:"not null"`
	Experience     int       `json:"experience" gorm:"not null"` // In years
	ProfilePicture string    `json:"profile_picture" gorm:"not null"`
	Latitude       float64   `json:"latitude" gorm:"not null"`
	Longitude      float64   `json:"longitude" gorm:"not null"`
	VerifiedStatus bool      `gorm:"default:false"`
	AdharID        string    `json:"aadhar" gorm:"not null"`
	DOB            time.Time `json:"not null"`

	// Services offered by this provider will be stored in a separate table.
	Services []Service `gorm:"foreignKey:ProviderID"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
