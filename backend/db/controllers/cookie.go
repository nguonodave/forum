package controllers

import "net/http"

func GetCookie(r *http.Request, token string) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err.Error() == "ErrNoCookie" {
			return "",  err
		}
		return "",  err
	}

	return cookie.Value, nil
}