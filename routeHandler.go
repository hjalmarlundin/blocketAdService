package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func handleGetRequest(w http.ResponseWriter, r *http.Request, db *gorm.DB) {

	urlParams := r.URL.Query()
	fmt.Printf("Received GET with parameter: %v \n", urlParams)
	query := r.FormValue("sort_by")

	ads := searchAndQuery(query, db)
	fmt.Println("{}", ads)
	json.NewEncoder(w).Encode(ads)
}

func searchAndQuery(query string, db *gorm.DB) []blocketAd {
	var ads []blocketAd

	switch query {
	case "price.desc":
		fmt.Println("sort by price desc")
		db.Order("price desc").Find(&ads)
	case "price.asc":
		fmt.Println("sort by price asc")
		db.Order("price asc").Find(&ads)
	case "date.asc":
		fmt.Println("sort by date asc")
		db.Order("Created asc").Find(&ads)
	case "date.desc":
		fmt.Println("sort by date desc")
		db.Order("Created desc").Find(&ads)
	default:
		fmt.Printf("No or unknown sorting options provided: %v. Defaulting to none. \n", query)
		db.Find(&ads)
	}
	return ads
}

func handleDeleteRequest(w http.ResponseWriter, r *http.Request, db *gorm.DB) {

	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("Received DELETE for ID: %v \n", id)

	var ad blocketAd
	result := db.Where("Id = ?", id).Find(&ad)

	if result.Error == nil && result.RowsAffected != 0 {
		db.Delete(&ad)
		fmt.Fprintf(w, "Successfully Ad with ID: %v", id)
	} else {
		fmt.Printf("Could not find ad with id: %v", id)
		http.NotFound(w, r)
	}
}

func handlePostRequest(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var ad blocketAd
	err := json.Unmarshal(reqBody, &ad)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Received POST: \n Subject: %v. Body: %v, Email: %v, Price: %v", ad.Subject, ad.Body, ad.Email, ad.Price)

	createdResource := &blocketAd{Subject: ad.Subject, Body: ad.Body, Email: ad.Email, Price: ad.Price, Created: time.Now(), ID: uuid.NewString()}
	db.Create(createdResource)
	json.NewEncoder(w).Encode(createdResource)
}

//Preparing for common validation, auth etc.
func makeHandler(fn func(http.ResponseWriter, *http.Request, *gorm.DB), db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, db)
	}
}
