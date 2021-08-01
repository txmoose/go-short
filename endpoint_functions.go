package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RedirectToTargetURL(w http.ResponseWriter, r *http.Request) {

}

func ShowSlugDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var slug Slug
	DB.First(&slug, "slug = ? ", params["slug"])
	json.NewEncoder(w).Encode(slug)
}

func CreateNewSlug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slug Slug
	json.NewDecoder(r.Body).Decode(&slug)
	DB.Create(&slug)
	json.NewEncoder(w).Encode(slug)
}

func CreateCustomSlug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

func ShowRecentSlugs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slugs []Slug
	DB.Find(&slugs)
	json.NewEncoder(w).Encode(slugs)
}
