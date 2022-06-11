package main

import (
	"errors"
	"fmt"
	"log"

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
	// golang:GoTmpPasswd@tcp(zealot.lan:3306)/go-practice?parseTime=true
	// Defaults for DB Connection Info
	viper.SetDefault("config.db_user", "go-short")
	viper.SetDefault("config.db_host", "mysql")
	viper.SetDefault("config.db_port", "3306")
	viper.SetDefault("config.db_name", "go-short")
	err := viper.BindEnv("config.db_pass", "GS_DB_PASS")
	if err != nil {
		log.Fatal("No Database Password Set", err.Error())
	}

	viper.BindEnv("config.db_user", "GS_DB_USER")
	viper.BindEnv("config.db_host", "GS_DB_HOST")
	viper.BindEnv("config.db_port", "GS_DB_PORT")
	viper.BindEnv("config.db_name", "GS_DB_NAME")

	dbUser := viper.GetString("config.db_user")
	dbHost := viper.GetString("config.db_host")
	dbPort := viper.GetString("config.db_port")
	dbName := viper.GetString("config.db_name")
	dbPass := viper.GetString("config.db_pass")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Can not connect to DB")
	}
	err = DB.AutoMigrate(&Slug{})
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

// GetURLFromDB convenience function to get existing URL records from database
func GetURLFromDb(url string) (Slug, error) {
	var slug Slug
	result := DB.Where("target_url = ? ", url).First(&slug)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return slug, result.Error
	}
	return slug, nil
}
