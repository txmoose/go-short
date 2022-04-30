package main

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Slug DB Table Mapping
type Slug struct {
	SiteTitle string `json:"site_title"`
	TargetURL string `json:"target_url"`
	Slug      string `gorm:"PrimaryKey"`
	HitCount  uint32 `json:"hit_count"`
}

// DB err Global variables for namespacing ease
var DB *gorm.DB
var err error

// InitializeDB creates database and does migrations if necessary
func InitializeDB() {
	// dsn database Connection string
	// TODO Is MySQL the best option?
	dsn := viper.GetString("config.dsn")
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
	result := DB.Select([]string{"slug"}).First(&slug, "slug = ? ", s)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return slug, result.Error
	}
	return slug, nil
}

// CheckTargetUrlExists returns true if a URL is already in the database
func CheckTargetUrlExists(url string) bool {
	var slug Slug
	result := DB.Select([]string{"target_url"}).First(&slug, "target_url = ? ", url)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
