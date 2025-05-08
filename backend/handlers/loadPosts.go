package handlers

import (
	"encoding/json"
	"net/http"
	"real_time_forum/backend/db/controllers"
)

func LoadAll(w http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("token")
	 if !controllers.ValidSession(token){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{"message" : "Please login to continue"})
		return
	 }

	all, err := controllers.All()

	if err != nil {
		// log.Println(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"posts" : err.Error()})
		return
	}

	if len(all) == 0 {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"posts": []interface{}{}})
	return
}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"posts" : all})
}