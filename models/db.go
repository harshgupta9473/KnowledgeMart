package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Email    string `gorm:"type:varchar(255);unique" validate:"required,email" json:"email" `
	Password string `gorm:"type:varchar(255)" validate:"required" json:"password"`
}

type User struct {
	gorm.Model
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string  `gorm:"type:varchar(255)" validate:"required" json:"name"`
	Email        string  `gorm:"type:varchar(255);unique" validate:"email" json:"email"`
	PhoneNumber  string  `gorm:"type:varchar(255);unique" validate:"number" json:"phone_number"`
	Picture      string  `gorm:"type:text" json:"picture"`
	WalletAmount float64 `gorm:"type:double precision" json:"wallet_amount"`
	Pasword      string  `gorm:"type:varchar(255)" validate:"required" json:"password"`
	Blocked      bool    `gorm:"type:bool" json:"blocked"`
	ReferralCode string  `json:"referral_code"`
	OTP          string  `gorm:"type:varchar(40)"`
	OTPExpiry    time.Time
	IsVerified   bool   `gorm:"type:bool" json:"verified"`
	LoginMethod  string `gorm:"type:varchar(50)" json:"login_method"`
}

//  type UserReferralHistory struct

type Seller struct {
	gorm.Model
	ID            uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint    `gorm:"not null;constraint:OnDelete:CASCADE;" json:"userId"`
	UserName      string  `gorm:"type:varchar(255)" validate:"required" json:"name"`
	WalletAmount  float64 `gorm:"type:double precision" json:"wallet_amount"`
	Password      string  `gorm:"type:varchar(255)" validate:"required" json:"password"`
	Description   string  `gorm:"type:varchar(255)" validate:"required" json:"description"`
	IsVerified    bool    `gorm:"type:bool" json:"verified"`
	AverageRating float64 `gorm:"type:decimal(10,2)" json:"averageRating"`

	User User `gorm:"foreignKey:UserID"`
}//

type Category struct {
	gorm.Model
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string    `gorm:"type:varchar(255)" validate:"required" json:"name"`
	Description     string    `gorm:"type:varchar(255)" validate:"required" json:"image"`
	OfferPercentage uint      `json:"offer_percentage"`
	Image           string    `gorm:"type:varchar(255)" validate:"required" json:"image"`
	Products        []Product `gorm:"foreignKey:CategoryID"`
}
type Product struct {
	gorm.Model
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	SellerID     uint           `gorm:"not null;constraint:OnDelete:CASCADE;" json:"sellerId"`
	Name         string         `gorm:"type:varchar(255)" validate:"required" json:"name"`
	CategoryID   uint           `gorm:"not null; constraint:OnDelete:CASCADE;" json:"categoryId"`
	Description  string         `gorm:"type:varchar(255)" validate:"required" json:"description"`
	Availability bool         `gorm:"type:bool;default:true" json:"availabilty"`
	Price        float64        `gorm:"type:decimal(10,2);not null" validate:"required" json:"price"`
	OfferAmount  float64        `gorm:"type:decimal(10,2);not null" validate:"required" json:"offer_amount"`
	Image        pq.StringArray `gorm:"type:varchar(255)[]" validate:"required" json:"image_url"`

	Seller   Seller   `gorm:"foreignKey:SellerID"`
	Category Category `gorm:"foreignKey:CategoryID"`
}

type Address struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint   `gorm:"not null;constraint:OnDelete:CASCADE;" json:"userId"`
	StreetName   string `gorm:"type:varchar(255)" validate:"required" json:"street_name"`
	StreetNumber string `gorm:"type:varchar(255)" validate:"required" json:"street_number"`
	City         string `gorm:"type:varchar(255)" validate:"required" json:"city"`
	State        string `gorm:"type:varchar(255)" validate:"required" json:"state"`
	PinCode      string `gorm:"type:varchar(255)" validate:"required" json:"pincode"`
	PhoneNumber  string `gorm:"type:varchar(255);unique" validate:"number" json:"phone_number"`

	User         User   `gorm:"foreignKey:UserID"`
}

type Cart struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint    `gorm:"not null" json:"userId"`
	ProductID uint    `gorm:"not null" json:"productId"`

	User      User    `gorm:"foreignKey:UserID"`
	Product   Product `gorm:"foreignKey:ProductID"`
}

type Order struct {
	OrderID                uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                 uint            `gorm:"not null" json:"user_id"`
	CouponCode             string          `json:"coupon_code"`
	CouponDiscountAmount   float64         `validate:"required,number" json:"coupon_discount_amount"`
	CategoryDiscountAmount float64         `validate:"required,number" json:"category_discount_amount"`
	ProductOfferAmount     float64         `validate:"required,number" json:"product_offer_amount"`
	DeliveryCharge         float64         `validate:"number" json:"delivery_charge"`
	TotalAmount            float64         `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	FinalAmount            float64         `validate:"required,number" json:"final_amount"`
	PaymentMethod          string          `gorm:"type:varchar(100)" validate:"required" json:"payment_method"`
	PaymentStatus          string          `gorm:"type:varchar(100)" validate:"required" json:"payment_status"`
	OrderedAt              time.Time       `gorm:"autoCreateTime" json:"ordered_at"`
	ShippingAddress        ShippingAddress `gorm:"embedded" json:"shipping_address"`
	SellerID               uint            `gorm:"not null" json:"seller_id"`
	Status                 string          `gorm:"type:varchar(100);default:'pending'" json:"status"`
	FailedPaymentCount     int             `gorm:"default:0" json:"failed_Payment_count"`
}

type ShippingAddress struct {
	StreetName   string `gorm:"type:varchar(255)" json:"street_name"`
	StreetNumber string `gorm:"type:varchar(255)" json:"street_number"`
	City         string `gorm:"type:varchar(255)" json:"city"`
	State        string `gorm:"type:varchar(255)" json:"state"`
	PinCode      string `gorm:"type:varchar(20)" json:"pincode"`
	PhoneNumber  string `gorm:"type:varchar(20)" json:"phonenumber"`
}

type OrderItem struct {
	
	OrderItemID uint `gorm:"primaryKey;autoIncrement" json:"orderItemId"`
	OrderID     uint `gorm:"not null" json:"orderId"`
	UserID              uint    `gorm:"not null" json:"userId"`
	ProductID           uint    `gorm:"not null" json:"productId"`
	SellerID            uint    `gorm:"not null" json:"sellerId"`
	Price               float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	ProductOfferAmount  float64 `json:"product_offer_amount"`
	CategoryOfferAmount float64 `json:"category_offer_amount"`
	OtherOffers         float64 `json:"other_offers"`
	FinalAmount         float64 `json:"final_amount"`
	Status              string  `gorm:"type:varchar(100);default:'pending'" json:"status"`

	User                User    `gorm:"foreignKey:UserID"`
	Product             Product `gorm:"foreignKey:ProductID"`
	Seller              Seller  `gorm:"foreignKey:SellerID"`
}