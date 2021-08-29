package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/url"
)

// TODO figure out how to get this out of env vars
const SlugLength = 8

// initializeRouter the routes
func initializeRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/create", CreateNewSlug).Methods("POST")
	router.HandleFunc("/recent", ShowRecentSlugs).Methods("GET")
	router.HandleFunc("/{slug}", RedirectToTargetURL).Methods("GET")
	router.HandleFunc("/{slug}/detail", ShowSlugDetail).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

// RedirectToTargetURL what it says on the tin
func RedirectToTargetURL(w http.ResponseWriter, r *http.Request) {
	// This is how we get variables out of the URI with mux
	params := mux.Vars(r)
	var slug Slug
	slug, err = GetSlugFromDB(params["slug"])
	if err != nil {
		panic(err)
	}
	slug.HitCount++
	DB.Save(&slug)
	http.Redirect(w, r, slug.TargetURL, 301)
}

// ShowSlugDetail shows the details of a slug, so you can know your redirect is safe
func ShowSlugDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var slug Slug
	slug, err = GetSlugFromDB(params["slug"])
	// if GetSlugFromDB throws an error, we're gonna throw HTTP 500
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// See this line a lot, this is a JSON encoder writing to the Response Writer
	// To send a JSON back to the user
	err := json.NewEncoder(w).Encode(slug)
	// if JSON encoding fails, we throw an HTTP 500
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateNewSlug creates a new slug, will generate a random slug if necessary or use passed custom slug
func CreateNewSlug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slug Slug
	// Decode request body into a slug Struct
	err := json.NewDecoder(r.Body).Decode(&slug)
	// if JSON decoding fails, we throw an HTTP 400
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Parse the URL in the body, and if it is invalid, tell the user
	u, err := url.Parse(slug.TargetURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// if slug is not passed as part of the request body, we generate a random one
	if slug.Slug == "" {
		slug.Slug, err = GenerateSlug(SlugLength)
		// if we have a generation error, throw HTTP 500
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// if slug is passed as part of the request body, we ensure it doesn't already exist
		//TODO implement common word list to also disallow
	} else {
		slug, err = GetSlugFromDB(slug.Slug)
		// this is confusing, but if we get Record Not Found Error, we're good to continue
		// but if we get anything _other_ than Record Not Found, we throw HTTP 400 and let user know
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Slug already in Use", http.StatusBadRequest)
		}
	}

	// Go out and get the title of the site
	slug.SiteTitle, err = GetSiteTitle(slug.TargetURL)
	if err != nil {
		slug.SiteTitle = u.Hostname()
	}

	DB.Create(&slug)
	err = json.NewEncoder(w).Encode(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//ShowRecentSlugs should show only N most recent slugs
//TODO limit this to some number N
func ShowRecentSlugs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slugs []Slug
	DB.Find(&slugs)
	err := json.NewEncoder(w).Encode(slugs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	InitializeDB()
	initializeRouter()
}
