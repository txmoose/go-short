package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Slug struct {
	gorm.Model
	SiteTitle string `json:"site_title"`
	TargetURL string `json:"target_url"`
	Slug      string `json:"slug"`
	HitCount  uint32 `json:"hit_count"`
}

var DB *gorm.DB
var err error

const dsn = "golang:password!@tcp(zealot.lan:3306)/go-practice"

func InitializeDB() {
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Can not connect to DB")
	}
	DB.AutoMigrate(&Slug{})
}

func RedirectToTargetURL(w http.ResponseWriter, r *http.Request) {

}

func ShowSlugDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var slug Slug
	DB.Select([]string{"SiteTitle", "TargetURL", "Slug", "HitCount"}).Find(&slug, "slug = ?", params["id"])
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
	result := DB.Find(&slugs)
	json.NewEncoder(w).Encode(result.Rows)
}

func initializeRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/{slug}", RedirectToTargetURL).Methods("GET")
	router.HandleFunc("/{id}/detail", ShowSlugDetail).Methods("GET")
	router.HandleFunc("/create", CreateNewSlug).Methods("POST")
	router.HandleFunc("/custom", CreateCustomSlug).Methods("POST")
	router.HandleFunc("/recent", ShowRecentSlugs).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	InitializeDB()
	initializeRouter()
}
