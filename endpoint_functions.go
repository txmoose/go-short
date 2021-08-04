package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"log"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"
)

const SlugLength = 8

func initializeRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/create", CreateNewSlug).Methods("POST")
	router.HandleFunc("/recent", ShowRecentSlugs).Methods("GET")
	router.HandleFunc("/{slug}", RedirectToTargetURL).Methods("GET")
	router.HandleFunc("/{slug}/detail", ShowSlugDetail).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func RedirectToTargetURL(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var slug Slug
	slug, err = GetSlugFromDB(params["slug"])
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, slug.TargetURL, 301)
}

func ShowSlugDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var slug Slug
	slug, err = GetSlugFromDB(params["slug"])
	if err != nil {
		panic(err)
	}
	err := json.NewEncoder(w).Encode(slug)
	if err != nil {
		return
	}
}

func CreateNewSlug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slug Slug
	err := json.NewDecoder(r.Body).Decode(&slug)
	if err != nil {
		return
	}
	if slug.Slug == "" {
		slug.Slug, err = GenerateSlug(SlugLength)
		if err != nil {
			panic(err)
		}
	}

	DB.Create(&slug)
	err = json.NewEncoder(w).Encode(slug)
	if err != nil {
		return
	}
}

func ShowRecentSlugs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var slugs []Slug
	DB.Find(&slugs)
	err := json.NewEncoder(w).Encode(slugs)
	if err != nil {
		return
	}
}

// GenerateSlug Found and modified from here
// https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb#gistcomment-3527095
func GenerateSlug(n int) (string, error) {
	//goland:noinspection ALL
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	var slug Slug
	slug, _ = GetSlugFromDB(string(ret))
	candidateSlug := string(ret)

	if candidateSlug == slug.Slug {
		candidateSlug, err = GenerateSlug(n)
	}
	return candidateSlug, nil
}

func GetSlugFromDB(s string) (Slug, error) {
	var slug Slug
	result := DB.First(&slug, "slug = ? ", s)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return slug, result.Error
	}
	return slug, nil
}
