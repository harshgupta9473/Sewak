package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/harshgupta9473/sevak_backend/config"
	"github.com/harshgupta9473/sevak_backend/middleware"
	"github.com/harshgupta9473/sevak_backend/models"
	"github.com/harshgupta9473/sevak_backend/utils"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{
		db: db,
	}
}

func (auth *AuthController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if err := godotenv.Load(); err != nil {
		http.Error(w, "server error in loading .env", http.StatusInternalServerError)
		return
	}
	lenString := os.Getenv("lengthOfOTP")
	length, err := strconv.Atoi(lenString)
	if err != nil {
		http.Error(w, "server error in converting length to int", http.StatusInternalServerError)
		return
	}

	var userReq models.RUser
	err = json.NewDecoder(r.Body).Decode(&userReq)

	if !config.IfRoleExistINOURSYSTEM(userReq.Role) {
		http.Error(w, "not allowed", http.StatusForbidden)
		return
	}

	if err != nil {
		http.Error(w, "invalid formate", http.StatusBadRequest)
		return
	}

	var user models.User
	result := auth.db.Preload("Roles").Where("mobile = ?", userReq.Mobile).First(&user)
	if result.Error != nil {
		// if user does not exists
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// user does not exists yet
			// is role is admin
			if utils.IsAdmin(userReq.Role) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			} else {
				// create account
				err = auth.Account(&userReq, length)
				if err != nil {
					http.Error(w, "error occured", http.StatusInternalServerError)
					return
				}
				utils.WriteJSON(w, http.StatusOK, "success")
				return
			}
		} else {
			// some internal datbase connection error may be
			http.Error(w, "error occured", http.StatusInternalServerError)
			return
		}
	} else {
		// user does exists

		// checking if role provided is actually alloted to user or not
		if utils.IsRoleAlloted(user.Roles, userReq.Role) {

			// login
			err = auth.Account(&userReq, length)
			if err != nil {
				http.Error(w, "error occured", http.StatusInternalServerError)
				return
			}
			utils.WriteJSON(w, http.StatusOK, "success")
			return

		} else {
			//role alloted does not have role that is provided different so create account

			// checking if roleprovided is admin? if admin as admin role is not allowed to created return forbidden
			if utils.IsAdmin(userReq.Role) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			} else {
				// role is not admin

				// create account
				err = auth.Account(&userReq, length)
				if err != nil {
					http.Error(w, "error occured", http.StatusInternalServerError)
					return
				}
				utils.WriteJSON(w, http.StatusOK, "success")
				return

			}
		}
	}

}

func (auth *AuthController) Account(userReq *models.RUser, length int) error {
	otp, err := utils.GenerateOTP(length)
	if err != nil {
		log.Println(err)
		return err
	}
	err = auth.SaveOTP(userReq.Mobile, otp, userReq.Role)
	if err != nil {
		log.Println(err)
		return err
	}
	err = utils.SendOtp(userReq.Mobile, otp)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (auth *AuthController) SaveOTP(mobile string, otp string, role string) error {
	var tempuser models.TempUser
	result := auth.db.Where("mobile=? AND role=?", mobile, role).First(&tempuser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err := auth.db.Create(&models.TempUser{Mobile: mobile, OTP: otp, Role: role, Time: time.Now().Add(5 * time.Minute)}).Error
			if err != nil {
				return err
			}
		} else {
			return result.Error
		}
	} else {
		tempuser.OTP = otp
		tempuser.Time = time.Now().Add(5 * time.Minute)
		err := auth.db.Save(tempuser).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (auth *AuthController) HandleVerification(w http.ResponseWriter, r *http.Request) {
	var usereq models.VUser
	err := json.NewDecoder(r.Body).Decode(&usereq)
	if err != nil {
		http.Error(w, "invalid format", http.StatusBadRequest)
		return
	}
	if !config.IfRoleExistINOURSYSTEM(usereq.Role) {
		http.Error(w, "not allowed", http.StatusForbidden)
		return
	}
	var TUser models.TempUser
	result := auth.db.Where("mobile=? AND role=?", usereq.Mobile, usereq.Role).First(&TUser)
	if result.Error != nil {
		http.Error(w, "invalid otp", http.StatusForbidden)
		return
	}
	if TUser.OTP != usereq.OTP {
		http.Error(w, "invalid otp", http.StatusForbidden)
		return
	}
	if time.Now().After(TUser.Time) {
		http.Error(w, "otp is expired", http.StatusGatewayTimeout)
		return
	}
	// user is authorised ok
	// check user already exists or not? if yes then log in him if no then insert him in database and login him
	var user models.User
	result = auth.db.Preload("Roles").Where("mobile=?", usereq.Mobile).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// create account
			user = models.User{Mobile: TUser.Mobile, Roles: []*models.Role{{Name: TUser.Role}}}
			err := auth.db.Save(&user)
			if err != nil {
				http.Error(w, "error occured while creating account ", http.StatusInternalServerError)
				return
			}
		} else {
			// db error
			http.Error(w, "error occured internal", http.StatusInternalServerError)
			return

		}
	} else {
		// found but let's check if  role  request matches the the existing role and if matches then it means login and if not then account creation
		if !utils.IsRoleAlloted(user.Roles, usereq.Role) {
			// create account
			user.Roles = append(user.Roles, &models.Role{Name: usereq.Role})
			err = auth.db.Save(user).Error
			if err != nil {
				http.Error(w, "error in account creation try after sometime", http.StatusInternalServerError)
				return
			}

		}
		// log in using cookies
	}

	var userFromDB models.User

	result = auth.db.Preload("Roles").Where("mobile=?", usereq.Mobile).First(&user)
	if result.Error!=nil{
		http.Error(w, "error in account creation try after sometime 2", http.StatusInternalServerError)
			return
	}

	loginToken, err := middleware.CreateLoginJWT(&userFromDB, TUser.Role)
	if err != nil {
		log.Println(err)
		http.Error(w, "error occured in login jwt", http.StatusInternalServerError)
		return
	}
	refreshtoken, err := middleware.CreateRefreshJWT(&userFromDB, TUser.Role)
	if err != nil {
		log.Println(err)
		http.Error(w, "error occured refrsh jwt", http.StatusInternalServerError)
		return
	}
	utils.SetAuthTokenCookie(w, loginToken, "authToken")
	utils.SetAuthTokenCookie(w, refreshtoken, "refreshToken")
	utils.WriteJSON(w, http.StatusOK, "success")

}
