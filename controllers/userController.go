package controllers

import (
	"fmt"
	"net/http"
)
func UserHandler(w http.ResponseWriter, r *http.Request) {
	//status http
	fmt.Fprintf(w, "que onda pa")
	return
}

func userRegister(username string, password string) bool {
	return false
}

func userLogin(username string, password string) bool {

	return true
}
