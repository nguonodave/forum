package controllers

import (
	"net/http"
)


func GetUserNAme(r *http.Request) ( string,string) {

	var name string
	cookie, _ := GetCookie(r, "session_token")
	
	var uname, uid string
	if ValidSession(cookie) {
		uid, uname, _, _ = GetUserDetailsFromSession(cookie)
		name = uname
	} else {
		name = "GUEST"
	}

	return uid, name
}