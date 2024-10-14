package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/harshgupta9473/sevak_backend/models"
	"github.com/harshgupta9473/sevak_backend/utils"
	"gorm.io/gorm"
)

type ProfileController struct {
	db *gorm.DB
}

func NewProfileController(db *gorm.DB) *ProfileController {
	return &ProfileController{
		db: db,
	}
}

func (pc *ProfileController) CreateOrUpdateCustomerProfileByCustomer(w http.ResponseWriter, r *http.Request) {
	var profile models.CustomerProfile
	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if !utils.IsRoleMatches("customer",r){
		http.Error(w,"unauthorised",http.StatusOK)
		return
	}
	
	result := pc.db.Save(&profile)
	if result.Error != nil {
		http.Error(w, "Failed to save profile", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w,http.StatusOK,profile)
}


func (pc *ProfileController) CreateOrUpdateServiceProviderProfileByServiceProvider(w http.ResponseWriter, r *http.Request) {
    var profile models.ServiceProviderProfile
    err := json.NewDecoder(r.Body).Decode(&profile)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
	if !utils.IsRoleMatches("worker",r){
		http.Error(w,"unauthorised",http.StatusOK)
		return
	}
	var existingProfile models.ServiceProviderProfile
	result:=pc.db.First(&existingProfile,"user_id",profile.UserID)
	if result.Error!=nil && !errors.Is(result.Error,gorm.ErrRecordNotFound){
		http.Error(w,"failed to update",http.StatusInternalServerError)
		return
	}
    
	if result.Error!=nil{
		profile.VerifiedStatus=false
	}else{
		profile.VerifiedStatus=existingProfile.VerifiedStatus
	}
	result = pc.db.Save(&profile)
    if result.Error != nil {
        http.Error(w, "Failed to save profile", http.StatusInternalServerError)
        return
    }
    utils.WriteJSON(w,http.StatusOK,profile)
}

