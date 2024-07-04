package db

import (
	"log"
	"paysyncevets/utils"
	"sync"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"

)

var (
	db *gorm.DB //database
	err error
	once sync.Once
)

func InitDB(config utils.Config) {
	once.Do(func() {
		var error error
        db, err = gorm.Open(postgres.Open(config.DBSource), &gorm.Config{})
		if error != nil {
			panic(error)
		}
		log.Println("Database connection established")
	})
	
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database connection not established")
	}
	return db
}