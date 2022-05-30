package models

import (
	"fmt"
	"log"

	"github.com/2022AA/bytes-linked/backend/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

// Setup initializes the database instance
func dbSetUp(uri string, dbType string) {
	var err error
	db, err = gorm.Open(dbType, uri)

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}
func Setup() {
	uri := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name)
	dbSetUp(uri, setting.DatabaseSetting.Type)
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	defer db.Close()
}
