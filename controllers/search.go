package controllers

import (
	"net/http"

	database "github.com/harshgupta9473/NotesCart/config"
	"github.com/harshgupta9473/NotesCart/models"
	"github.com/harshgupta9473/NotesCart/utils"
)

func SearchProducts(w http.ResponseWriter,r *http.Request){
	var products []models.Product
	var productResponse []models.ProductResponse

	recievedQuery:=r.URL.Query();
	categoryID:=recievedQuery.Get("category_id")
	sortBy:=recievedQuery.Get("sort_by")
	filterAvailable:=recievedQuery.Get("available")
	
	//
	query:=database.DB.Model(&products)

	if filterAvailable=="true"{
		query=query.Where("availability=?",true)
	}

	if categoryID!=""{
		query=query.Where("category_id=?",categoryID)
	}

	switch sortBy{
	case "price_asc":
		query=query.Order("offer_amount ASC")
	case "price_desc":
		query=query.Order("offer_amount DESC")
	case "newest":
		query=query.Order("created_at DESC")	
	case "name_asc":
		query=query.Order("LOWER(name) DESC")
	case "high_rating":
		query=query.Joins("JOIN sellers ON sellers.id=products.seller_id").Order("sellers.average_rating DESC")			
	}

	tx:=query.Find(&products)

	if tx.Error!=nil{
		utils.WriteJSON(w,http.StatusNotFound,map[string]string{
			"status":"failed",
			"message":"failed to retrieve data from the products database, or the data doesn't exists",
		})
		return 
	}

	for _,product:=range products{

		var seller models.Seller
		if err:=database.DB.Where("id=?",product.SellerID).First(&seller).Error;err!=nil{
			utils.WriteJSON(w,http.StatusInternalServerError,map[string]string{
				"status":"failed",
				"message":"failed to retrieve seller",
			})
			return 
		}
		productResponse=append(productResponse, models.ProductResponse{
			ID: product.ID,
			Name: product.Name,
			Description: product.Description,
			Price: product.Price,
			OfferAmount: product.OfferAmount,
			Image: product.Image,
			Availability: product.Availability,
			SellerID: product.SellerID,
			SellerName: seller.UserName,
			CategoryID: product.CategoryID,
			SellerRating: seller.AverageRating,
		})
	}

	utils.WriteJSON(w,http.StatusOK,map[string]interface{}{
		"status":"success",
		"messgae":"successfully retrieved products",
		"data":map[string]interface{}{
			"products":productResponse,
		},
	})

}


