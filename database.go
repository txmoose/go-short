package main

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Slug DB Table Mapping
type Slug struct {
	ID        uint   `gorm:"primaryKey"`
	SiteTitle string `json:"site_title"`
	TargetURL string `json:"target_url"`
	Slug      string `gorm:"index,unique"`
	HitCount  uint32 `json:"hit_count"`
}

// DB err Global variables for namespacing ease
var DB *gorm.DB
var err error

// dsn database Connection string
// TODO Make this pull from env vars rather than hard code
// TODO Is MySQL the best option?
const dsn = "golang:password!@tcp(zealot.lan:3306)/go-practice?parseTime=true"

// InitializeDB creates database and does migrations if necessary
func InitializeDB() {
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Can not connect to DB")
	}
	err := DB.AutoMigrate(&Slug{})
	if err != nil {
		return
	}
}

// GetSlugFromDB convenience function to get slugs from database
func GetSlugFromDB(s string) (Slug, error) {
	var slug Slug
	// This omits the ID, but it still displays as a 0 rather than the actual ID
	// TODO Figure out how to omit ID field from the JSON completely on select
	result := DB.Select([]string{"site_title", "target_url", "slug", "hit_count"}).First(&slug, "slug = ? ", s)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return slug, result.Error
	}
	return slug, nil
}
