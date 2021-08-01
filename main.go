package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func initializeRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/create", CreateNewSlug).Methods("POST")
	router.HandleFunc("/custom", CreateCustomSlug).Methods("POST")
	router.HandleFunc("/recent", ShowRecentSlugs).Methods("GET")
	router.HandleFunc("/{slug}", RedirectToTargetURL).Methods("GET")
	router.HandleFunc("/{slug}/detail", ShowSlugDetail).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	InitializeDB()
	initializeRouter()
}
