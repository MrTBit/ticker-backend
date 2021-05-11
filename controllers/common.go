package controllers

import (
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

func GetDb(w http.ResponseWriter, r *http.Request) (*gorm.DB, bool) {
	db, ok := r.Context().Value("DB").(*gorm.DB)

	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, false
	}

	return db, true
}

func ToJson(v interface{}, w http.ResponseWriter) ([]byte, bool) {
	jsonBytes, err := json.Marshal(v)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, false
	}

	return jsonBytes, true
}
