package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Initialize Viper for config and bail out if things are missing
func init() {
	// Config file name and location
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Setting prefix for ENV VARS
	viper.SetEnvPrefix("gs")

	// Checking for the SlugLength ENV VAR
	viper.BindEnv("config.slug_length", "GS_SLUG_LENGTH")

	// Some debug output
	log.Printf("Max Slug Length: %d", viper.GetInt("config.slug_length"))
}

// initializeRouter the routes
func initializeRouter() {
	log.Print("Initializing Router on port :8000")
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
	log.Printf("Redirecting %s to %s", slug.Slug, slug.TargetURL)
	http.Redirect(w, r, slug.TargetURL, 301)
	return
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
		return
	}
	// See this line a lot, this is a JSON encoder writing to the Response Writer
	// To send a JSON back to the user
	log.Printf("Getting Details for %s", slug.Slug)
	err := json.NewEncoder(w).Encode(slug)
	// if JSON encoding fails, we throw an HTTP 500
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateNewSlug creates a new slug, will generate a random slug if necessary or use passed custom slug
func CreateNewSlug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slug Slug
	SlugLength := viper.GetInt("config.slug_length")

	// Decode request body into a slug Struct
	err := json.NewDecoder(r.Body).Decode(&slug)

	// if JSON decoding fails, we throw an HTTP 400
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the URL in the body, and if it is invalid, tell the user
	u, err := url.Parse(slug.TargetURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	// if slug is not passed as part of the request body, we generate a random one
	if slug.Slug == "" {
		slug.Slug, err = GenerateSlug(SlugLength)
		// if we have a generation error, throw HTTP 500
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// if slug is passed as part of the request body, we ensure it doesn't already exist
		//TODO implement common word list to also disallow
	} else {
		_, err = GetSlugFromDB(slug.Slug)
		// this is confusing, but if we get Record Not Found Error, we're good to continue
		// but if we get anything _other_ than Record Not Found, we throw HTTP 400 and let user know
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Slug already in Use", http.StatusConflict)
			return
		}
	}

	existingUrl, err := GetURLFromDb(slug.TargetURL)
	if err == nil {
		log.Printf("Existing record, returning %s", existingUrl.Slug)
		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(existingUrl)
		// if JSON encoding fails, we throw an HTTP 500
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Go out and get the title of the site
	log.Printf("getting a site title for %s", slug.TargetURL)
	slug.SiteTitle, err = GetSiteTitle(u.String())
	if err != nil {
		log.Printf("Didn't get a site title for %s", slug.TargetURL)
		slug.SiteTitle = u.Hostname()
		log.Printf("Using %s", slug.SiteTitle)
	}

	DB.Create(&slug)
	err = json.NewEncoder(w).Encode(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//ShowRecentSlugs should show only N most recent slugs
func ShowRecentSlugs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slugs []Slug
	log.Print("Recent Slugs")
	DB.Limit(10).Find(&slugs)
	err := json.NewEncoder(w).Encode(slugs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	InitializeDB()
	initializeRouter()
}
