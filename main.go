package main

import (
	"log"
	"paysyncevets/api"
	"paysyncevets/db"
	"paysyncevets/models"
	"paysyncevets/utils"
)


func main() {

	// load the config file

	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// connect to postgress database

	db.InitDB(config)

	db := db.GetDB()

	err = db.AutoMigrate(
        &models.User{},
        &models.UserRole{},
        &models.Promoter{},
        &models.Artist{},
        &models.Venue{},
        &models.EventVenue{},
        &models.Booking{},
        &models.Ticket{},
		&models.Event{},
		&models.EventArtist{},
    )
    if err != nil {
        log.Fatal("failed to migrate models:", err)
    }

    log.Println("Database connected and models migrated successfully.")


	// create a server
	server, err := api.NewServer(config, db)

	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// start the server
	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}