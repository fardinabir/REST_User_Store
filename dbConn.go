package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB connection params
const (
	username = "root"
	password = "mypass98"
	hostname = "127.0.0.1:3306"
	dbname   = "users_db"
	extra    = "charset=utf8mb4&parseTime=true"
)

var db *sql.DB
var alreadyInitialized bool
var errorInDB bool

var DB *gorm.DB
var err error
var DSN string

func intialMigration() {
	dsn()
	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("cannot connect DB")
	}
	DB.AutoMigrate(&User{}, &UserTag{})
}

func dsn() {
	DSN = fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", username, password, hostname, dbname, extra)
}
