package main

import (
	"github.com/SongCastle/KoR/api/router"
	"github.com/SongCastle/KoR/internal/log"
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
		log.Fatalf("Failed to launch Server")
		return
	}
	router.Routes().Run(":8080")
}
