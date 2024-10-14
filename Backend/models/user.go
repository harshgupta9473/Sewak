package models

import (
	"time"

	"gorm.io/gorm"
)




type RUser struct{
	Mobile string  `json:"number"`
	Role   string  `json:"role"`
}

type VUser struct{
	Mobile string `json:"number"`
	OTP    string  `json:"otp"`
	Role   string  `json:"role"`
}



type TempUser struct {
    ID     int       `gorm:"primaryKey;autoIncrement" json:"id"`
    Mobile string    `gorm:"not null" json:"mobile"`
    OTP    string    `gorm:"not null" json:"otp"`
    Role   string    `gorm:"not null" json:"role"`
    Time   time.Time `gorm:"not null" json:"time"`
}

type Role struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

type User struct {
	gorm.Model
	Mobile string  `gorm:"unique"`
	Roles  []*Role `gorm:"many2many:user_roles"`
}

type Customer struct {
	gorm.Model
	UserID    uint `gorm:"not null"`
	Mobile    string
	FirstName string
	LastName  string
	CurrentLat     float64
	CurrentLong    float64
	Preferences    string // JSON or comma-separated values for preferences
	Address        string // Example: User's home address
	AdditionalInfo string // Any additional information

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Worker struct {
	gorm.Model
	UserID             uint `gorm:"not null"` // Foreign key to User
	Mobile             string
	Name               string `gorm:"not null"`
	Adhar              string `gorm:"unique;not null"`
	DOB                string `gorm:"unique; not null"`
	Address            string
	VerificationStatus string
	AdditionalDetails  string

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// type IndividualProvider struct {
// 	ServiceProvider
// 	Adhar string `gorm:"unique;not null"`
// }

// type CompanyProvider struct {
// 	ServiceProvider
// 	GSTNumber string `gorm:"unique;not null"` // GST number for services in India
// }

func InitUser(db *gorm.DB) {

	db.AutoMigrate(&User{}, &Role{}, &TempUser{},&Customer{},&Worker{},&CustomerProfile{},&ServiceProviderProfile{})

	adminRole := Role{Name: "Admin"}
	providerRole := Role{Name: "worker"}
	customerRole := Role{Name: "Customer"}
	
	db.FirstOrCreate(&adminRole)
	db.FirstOrCreate(&providerRole)
	db.FirstOrCreate(&customerRole)
}

