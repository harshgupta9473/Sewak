package utils

import (

	"github.com/harshgupta9473/sevak_backend/database"
	"github.com/harshgupta9473/sevak_backend/models"
)

func AccessUSERSByMobileNumber(mobile string) (*models.User,error){
	db:=database.GetDB()
	var user models.User
	result:=db.Preload("Roles").Where(`mobile=`,mobile).First(&user)
	if result.Error!=nil{
		return nil,result.Error
	}
	return &user,nil
}

