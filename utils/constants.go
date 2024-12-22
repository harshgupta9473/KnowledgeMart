package utils


const (
    RoleUser  = "user"
    RoleSeller = "seller"
    RoleAdmin = "admin"
)

type contextKey string

const (
	UserID contextKey = "userID"
	SellerID contextKey ="sellerID"
	AdminID contextKey ="adminID"
)