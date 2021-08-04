package main

import (
	"fmt"

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

const dsn = "golang:password!@tcp(zealot.lan:3306)/go-practice?parseTime=true"

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
