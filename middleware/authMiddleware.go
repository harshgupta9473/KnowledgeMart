package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/harshgupta9473/NotesCart/utils"
)



func Authenticate(next http.Handler)http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
        tokenstring:=w.Header().Get("Authorization")

    if tokenstring==""{
        utils.WriteJSON(w,http.StatusUnauthorized,map[string]string{
            "status":  "failed",
			"message": "Authorization header required",
        })
        return
    }
    token:=strings.TrimSpace(strings.Replace(tokenstring,"Bearer ", "", 1))
    if token==""{
        utils.WriteJSON(w,http.StatusUnauthorized,map[string]string{
            "status":  "failed",
			"message": "Invalid token format",
        })
        return
    }
    claims,err:=utils.ValidateJWT(token)
    if err!=nil{
        utils.WriteJSON(w,http.StatusUnauthorized,map[string]string{
            "status":  "failed",
			"message": fmt.Sprintf("Token validation failed:%v",err),
        })
        return
    }
    switch claims.Role{
    case "user":
       r= r.WithContext(context.WithValue(r.Context(),utils.UserID,claims.ID))
    case "seller":
       r= r.WithContext(context.WithValue(r.Context(),utils.SellerID,claims.ID))
    case "admin":
       r= r.WithContext(context.WithValue(r.Context(),utils.AdminID,claims.ID))
    default:
        utils.WriteJSON(w,http.StatusUnauthorized,map[string]string{
            "status":  "failed",
			"message": "unauthorized role",
        })
        return
    }
        next.ServeHTTP(w,r)
    })
	
}