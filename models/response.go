package models

import "github.com/lib/pq"

type UserProfileResponse struct {
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	PhoneNumber  string  `json:"phone_number"`
	Picture      string  `json:"picture"`
	ReferralCode string  `json:"referral_code"`
	WalletAmount float64 `json:"wallet_amount"`
}

type ProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	OfferAmount float64 `json:"offer_amount"`
	Image       pq.StringArray `json:"image_url"`
	Availability bool `json:"availability"`
	SellerID uint `json:"sellerid"`
	SellerName string `json:"sellername"`
	CategoryID uint `json:"categoryid"`
	SellerRating float64 `json:"sellerRating"`
}//

