package db

import (
	"log"
	"paysyncevets/models"
	"paysyncevets/utils"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB //database
	err error
	once sync.Once
)

func InitDB(config utils.Config) {
	once.Do(func() {
        db, err = gorm.Open(postgres.Open(config.DBSource), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		log.Println("Database connection established")
	})
	
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database connection not established")
	}
	seedRoles(db)
	return db
}

func seedRoles(db *gorm.DB) {
	roles := []models.UserRole{
		{RoleName: models.RoleAdmin},
		{RoleName: models.RoleArtist},
		{RoleName: models.RolePromoter},
		{RoleName: models.RoleNormal},
	}

	for _, role := range roles {
		// check if a role already exists
		var roleExists models.UserRole
		db.Where("role_name = ?", role.RoleName).First(&roleExists) 

		if roleExists.ID == 0 { // if the role does not exist
			db.Create(&role) // create the role
			log.Printf("Role %s created", role.RoleName)
		} else {
			log.Printf("Role %s already exists", role.RoleName)
		}
		
	}
}