package main

import "net/http"

func userHandler(w http.ResponseWriter, r *http.Request) {
	//status http
	w.WriteHeader(http.StatusOK)

	return
}

func userRegister(username string, password string) bool {
	return false
}

func userLogin(username string, password string) bool {

	return true
}
