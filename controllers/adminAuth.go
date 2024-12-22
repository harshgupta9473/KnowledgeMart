package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	database "github.com/harshgupta9473/NotesCart/config"
	"github.com/harshgupta9473/NotesCart/models"
	"github.com/harshgupta9473/NotesCart/utils"
	"gorm.io/gorm"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	err:=json.NewDecoder(r.Body).Decode(&loginData)
	if err!=nil{
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status":"failed",
			"error":"invalid request format",
		})
		return
	}
	var admin models.Admin
	if tx:=database.DB.Where("email=?",loginData.Email).First(&admin); tx.Error!=nil{
		if errors.Is(tx.Error,gorm.ErrRecordNotFound){
			w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":"failed",
				"error":"invalid email or password",
			})
			return
		}

		w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":"failed",
				"error":"database error",
			})
			return
	}
	if admin.Password!=loginData.Password{
		w.Header().Set("Content-Type","application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":"failed",
				"error":"invalid email or password",
			})
			return
	}
	token,err:=utils.CreateJWT(admin.ID,"admin")
	if token==""||err!=nil{
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status":"failed",
			"message":"failed to generate token",
		})
		return
	}
	w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token":token,
			"status":"sucess",
			"message":"login success",
		})
}
