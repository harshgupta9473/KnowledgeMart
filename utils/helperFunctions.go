package utils

import (
	"fmt"
	"net/http"
)

func GetIDFROMContext(r *http.Request,ctx contextKey)(error,uint) {
	id,ok:=r.Context().Value(ctx).(uint)
	if !ok{
		return fmt.Errorf("Unauthorized:Could not retrive user information"),0
	}
	return nil,id
}//