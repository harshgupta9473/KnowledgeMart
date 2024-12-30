package controllers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator"
	database "github.com/harshgupta9473/NotesCart/config"
	"github.com/harshgupta9473/NotesCart/models"
	"github.com/harshgupta9473/NotesCart/utils"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var Validate *validator.Validate

func EmailSignUp(w http.ResponseWriter, r *http.Request) {
	var signup models.EmailSignupRequest

	err := json.NewDecoder(r.Body).Decode(&signup)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"status":  "failed",
			"message": "failed to process the incoming request" + err.Error(),
		})
		return
	}
	Validate = validator.New()
	err = Validate.Struct(signup)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signup.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"status":  "failed",
			"message": "error in password hashing " + err.Error(),
		})
		return
	}
	refCode:=utils.GenerateRandomString(5)

	otp,otpExpiry,err:=generateAlphanumericOTP(6)
	if err!=nil{
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"status":  "failed",
			"message": "error in genrating otp " + err.Error(),
		})
		return
	}
	User:=models.User{
		Name: signup.Name,
		Email: signup.Email,
		PhoneNumber: signup.PhoneNumber,
		Blocked: false,
		Pasword: string(hashedPassword),
		ReferralCode: refCode,
		OTP: otp,
		OTPExpiry: otpExpiry,
		IsVerified: false,
		LoginMethod: "email",
	}

	tx:=database.DB.Where("email= ? AND deleted_at IS NULL",signup.Email).First(&User)
	if tx.Error!=nil && tx.Error!=gorm.ErrRecordNotFound{
		utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
			"status":  "failed",
			"message": "failed to retreive information from the database",
		})
		return
	}else if tx.Error==gorm.ErrRecordNotFound {
		tx=database.DB.Create(&User)
		if tx.Error!=nil{
			utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
				"status":  "failed",
				"message": "failed to create a new user",
			})
			return
		}
	}else{
		utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
			"status":  "failed",
			"message": "user already exist",
		})
		return
	}
	err=SendOTPMail(User.Email,otp)
	if err!=nil{
		utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
			"message": err.Error(),
		})
		return
	}
	// referral implementation is left



	//


	utils.WriteJSON(w,http.StatusOK,map[string]interface{}{
		"status":  "success",
		"message": "Account is Created Successfully, please login to complete your email verification",
		"data":map[string]interface{}{
			"name":         User.Name,
			"email":        User.Email,
			"phone_number": User.PhoneNumber,
			"picture":      User.Picture,
			"block_status": User.Blocked,
			"verified":     User.IsVerified,
		},
	})

}


func VerifyEmail(w http.ResponseWriter,r *http.Request){
	var emailOtpReq models.EmailVerificationRequest
	err:=json.NewDecoder(r.Body).Decode(&emailOtpReq)
	if  err!=nil{
		utils.WriteJSON(w,http.StatusBadRequest,map[string]string{
			"status":  "failed",
			"message": "wrong format",
		})
		return
	}
	emailOtpReq.OTP=strings.TrimSpace(emailOtpReq.OTP)

	validate:=validator.New()
	validate.Struct(emailOtpReq)

	var User models.User

	tx:=database.DB.Where("email=? and deleted_at IS NULL",emailOtpReq.Email).First(&User)
	if tx.Error!=nil{
		utils.WriteJSON(w,http.StatusBadRequest,map[string]string{
			"status":  "failed",
			"message": "user not found",
		})
		return
	}

	if User.OTP!=emailOtpReq.OTP || time.Now().After(User.OTPExpiry){
		utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
			"status":  "failed",
			"message": "invalid or expired otp",
		})
		return
	}

	User.IsVerified=true
	
	if err:=database.DB.Save(&User).Error; err!=nil{
		utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
			"status":  "failed",
			"message": "failed to update user",
		})
		return
	}

	utils.WriteJSON(w,http.StatusOK,map[string]string{
		"status":  "Success",
		"message": "Email Verified login with your credential",
	})
	
}


func EmailLogin(w http.ResponseWriter,r *http.Request){
	var LoginReq models.EmailLoginRequest
	err:=json.NewDecoder(r.Body).Decode(&LoginReq)
	if err!=nil{
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"status":  "failed",
			"message": "failed to process the incoming request" + err.Error(),
		})
		return
	}
	validate:=validator.New()
	err=validate.Struct(LoginReq)
	if err!=nil{
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"status":  "failed",
			"message": "" + err.Error(),
		})
		return
	}
	var User models.User
	tx:=database.DB.Where("email=? AND deleted_at IS NULL",LoginReq.Email).First(&User)
	if tx.Error!=nil{
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"status":  "failed",
			"message": "invalid email or password",
		})
		return
	}
	err=CheckPassword(User.Pasword,LoginReq.Password)
	if err!=nil{
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"status":  "failed",
			"message": "wrong email or password",
		})
		return
	}
	if User.Blocked{
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"status":  "failed",
			"message": "user is not authorised to access",
		})
		return
	}

	// /// // //
	if User.ReferralCode == "" {
		refCode := utils.GenerateRandomString(5)
		User.ReferralCode = refCode

		if err := database.DB.Model(&User).Update("referral_code", User.ReferralCode).Error; err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"status":  "failed",
				"message": "failed to update referral code",
			})
			return
		}
	}

	token,err:=utils.CreateJWT(User.ID,"user")
	if token==""||err!=nil{
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"status":  "failed",
			"message": "failed to generate token",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Login successful",
		"data":map[string]interface{}{
			"token":token,
			"id":           User.ID,
			"name":         User.Name,
			"email":        User.Email,
			"phone_number": User.PhoneNumber,
			"picture":      User.Picture,
			"block_status": User.Blocked,
			"verified":     User.IsVerified,
		},
	})
	
}

func ResendOTP(w http.ResponseWriter,r *http.Request){
	email:=r.URL.Query().Get("email")
	if email==""{
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"status":  "failed",
			"message": "email is required",
		})
		return
	}

	var user models.User
	tx:=database.DB.Where("email=? AND deleted_at IS NULL",email).First(&user)
	if tx.Error!=nil{
		if tx.Error==gorm.ErrRecordNotFound{
			utils.WriteJSON(w, http.StatusNotFound, map[string]string{
				"status":  "failed",
				"message": "user not found",
			})
			
		}else{
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
				"status":  "failed",
				"message": "failed to retrieve user",
			})
		}
		return
	}

	otp,otpExpiry,err:=generateAlphanumericOTP(6)
	if err!=nil{
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"status":  "failed",
			"message": "failed to generate otp",
		})
		return
	}
	user.OTP=otp
	user.OTPExpiry=otpExpiry

	tx=database.DB.Save(&user)
	if tx.Error!=nil{
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"status":  "failed",
			"message": "failed to update otp",
		})
		return
	}

	err=SendOTPMail(user.Email,otp)
	if err!=nil{
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"status":  "failed",
			"message": "failed to send OTP to mail",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "OTP has been sent successfully",
	})
	
}
//



//hashedPassword

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}



// otp genration

func generateAlphanumericOTP(length int) (string, time.Time, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	otp := make([]byte, length)
	for i := 0; i < length; i++ {
		max := big.NewInt(int64(len(charset)))
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", time.Time{}, err
		}
		otp[i] = charset[num.Int64()]
	}
	otpExpiry := time.Now().Add(3 * time.Minute)
	return string(otp), otpExpiry, nil
}

func generateNumericOTP(length int) (string, time.Time, error) {
	const charset = "0123456789"

	otp := make([]byte, length)
	for i := 0; i < length; i++ {
		max := big.NewInt(int64(len(charset)))
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", time.Time{}, err
		}
		otp[i] = charset[num.Int64()]
	}
	otpExpiry := time.Now().Add(3 * time.Minute)
	return string(otp), otpExpiry, nil
}



//sending otp via mail

func SendOTPMail(email, otp string) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	from := os.Getenv("emailID")
	password := os.Getenv("apppassword")
	smtpHost := os.Getenv("smtpHost")
	smtpPort := os.Getenv("smtpPort")

	msg := []byte("Subject: Verify your email\n\n" +
		fmt.Sprintf("Your OTP is %s", otp))

	auth := smtp.PlainAuth("", from, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(msg))

}
