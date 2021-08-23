package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) initialize(dbname string) {
	a.DB = createOrOpenDataBase(dbname)
	// Migrate the schema
	a.DB.AutoMigrate(&BlocketAd{})

	initializeRequestHandlers(a)

}

func createOrOpenDataBase(name string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func initializeRequestHandlers(app *App) {
	// creates a new instance of a mux router
	app.Router = mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	app.Router.HandleFunc("/ads", makeHandler(handleGetRequest, app.DB)).Methods("GET")
	app.Router.HandleFunc("/ad", makeHandler(handlePostRequest, app.DB)).Methods("POST")
	app.Router.HandleFunc("/ad/{id}", makeHandler(handleDeleteRequest, app.DB)).Methods("DELETE")
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
}

func (a *App) run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}