package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	database "github.com/harshgupta9473/NotesCart/config"
	"github.com/harshgupta9473/NotesCart/models"
	"github.com/harshgupta9473/NotesCart/utils"
	"gorm.io/gorm"
)

func SellerRegisteration(w http.ResponseWriter,r *http.Request) {
    var Register models.SellerRegisterRequest
    var user models.User
    var newSeller models.Seller

    err,userID:=utils.GetIDFROMContext(r,utils.UserID)
    if err!=nil{
        utils.WriteJSON(w,http.StatusUnauthorized,map[string]string{
            "status":"failed",
            "message":"user not authorised",
        })
        return
    }

    err=json.NewDecoder(r.Body).Decode(&Register)
    if err!=nil{
        utils.WriteJSON(w,http.StatusBadRequest,map[string]string{
            "status":"failed",
            "message":"failed to process the incoming request",
        })
        return
    }
    Validate=validator.New()

    err=Validate.Struct(Register)

    if err!=nil{
        utils.WriteJSON(w,http.StatusBadRequest,map[string]string{
            "status":"faiiled",
            "message":"missing fields",
        })
        return
    }
    if err:=database.DB.Where("id=? AND deleted_at IS NULL",userID).First(&user).Error;err!=nil{
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
    
    var seller models.Seller
    tx:=database.DB.Where("user_id=? AND deleted_at IS NULL",userID).First(&seller)
    if tx.Error==nil{
        utils.WriteJSON(w,http.StatusBadRequest,map[string]string{
            "status":"failed",
            "message":"seller already exists for this user",
        })
        return
    }else if tx.Error!=gorm.ErrRecordNotFound{
        utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
            "status":"failed",
            "message":"failed to retrieve seller information",
        })
        return
    }

    hashedPassword,err:=HashPassword(Register.Password)
    if err!=nil{
        utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
            "status":"failed",
            "message":"error in hashing password in seller registration",
        })
        return
    }

    newSeller=models.Seller{
        UserID: userID,
        UserName: Register.UserName,
        Password: hashedPassword,
        Description: Register.Description,
        IsVerified: false,
    }

    if err:=database.DB.Create(&newSeller).Error;err!=nil{
        utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
            "status":"failed",
            "message":"failed to create a new seller",
        })
        return
    }

    utils.WriteJSON(w,http.StatusOK,map[string]interface{}{
        "status":"success",
        "message":"Seller registered successfully, status is pending, please login to continue",
        "data":map[string]interface{}{
            "seller_name":user.Name,
            "seller_username":newSeller.UserName,
            "description":newSeller.Description,
            "email":user.Email,
            "phone_number":user.PhoneNumber,
            "user_id":newSeller.UserID,
            "sellerId":newSeller.ID,
        },
    })

}

func SellerLogin(w http.ResponseWriter,r *http.Request){
    var LoginSeller models.SellerLoginRequest

    err:=json.NewDecoder(r.Body).Decode(&LoginSeller)
    if err!=nil{
        utils.WriteJSON(w,http.StatusBadRequest,map[string]string{
            "status":"failed",
            "message":"unable to retrieve data from the seller login",
        })
        return
    }
    Validate=validator.New()
    err=Validate.Struct(LoginSeller)
    if err!=nil{
        utils.WriteJSON(w,http.StatusBadRequest,map[string]string{
            "status":"failed",
            "message":err.Error(),
        })
        return
    }

    var seller models.Seller
    tx:=database.DB.Where("user_id=? AND deleted_at is NULL",LoginSeller.UserID).First(&seller)
    if tx.Error!=nil{
        utils.WriteJSON(w,http.StatusBadRequest,map[string]string{
            "status":"failed",
            "message":"invalid email or password",
        })
        return
    }

    err=CheckPassword(seller.Password,LoginSeller.Password)
    if err!=nil{
        utils.WriteJSON(w,http.StatusUnauthorized,map[string]string{
            "status":"failed",
            "message":"incorrect email or password",
        })
        return
    }
    if !seller.IsVerified{
        utils.WriteJSON(w,http.StatusUnauthorized,map[string]string{
            "status":"failed",
            "message":"seller is not verified, status is pending",
        })
        return
    }

    token,err:=utils.CreateJWT(seller.ID,"seller")
    if token==""||err!=nil{
        utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
            "status":"failed",
            "message":"failed to generate token",
        })
        return
    }

    utils.WriteJSON(w,http.StatusOK,map[string]interface{}{
        "status":"success",
        "message":"Seller Login successful",
        "data":map[string]interface{}{
            "token":token,
            "username":seller.UserName,
            "verified":seller.IsVerified,
        },
    })
    

}
//