package controllers

import (
	"net/http"

	database "github.com/harshgupta9473/NotesCart/config"
	"github.com/harshgupta9473/NotesCart/models"
	"github.com/harshgupta9473/NotesCart/utils"
	"gorm.io/gorm"
)
//
func GetUserProfile(w http.ResponseWriter,r *http.Request){
	var user models.User

	err,userID:=utils.GetIDFROMContext(r,utils.UserID)
	if err!=nil{
		utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
			"status":"failed",
			"message":"failed to retrieve user information",
		})
		return
	}

	if err:=database.DB.Where("id=? AND deleted_at IS NULL",userID).First((&user)).Error;err!=nil{
		if err==gorm.ErrRecordNotFound{
			utils.WriteJSON(w,http.StatusBadRequest,map[string]string{
				"status":"failed",
				"message":"user not found",
			})
			return
		}

		utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
			"status":"failed",
			"message":"failed to retrieve user information",
		})
		return
	}

	userProfile:=models.UserProfileResponse{
		Name: user.Name,
		Email: user.Email,
		PhoneNumber: user.PhoneNumber,
		Picture: user.Picture,
		ReferralCode: user.ReferralCode,
		WalletAmount: user.WalletAmount,
	}
	utils.WriteJSON(w,http.StatusOK,map[string]interface{}{
		"status":  "success",
		"message": "successfully retrieved user profile",
		"data": map[string]interface{}{
			"profile": userProfile,
		},
	})

}

