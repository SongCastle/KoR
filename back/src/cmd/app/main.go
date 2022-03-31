package main

import (
	"log"

	"github.com/SongCastle/KoR/api/router"
	"github.com/SongCastle/KoR/volume/db"
)

func setUp() error {
	if err := db.SetUp(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := setUp(); err != nil {
		log.Printf("Failed to launch Server\n")
		return
	}
	router.Routes().Run(":8080")
}
