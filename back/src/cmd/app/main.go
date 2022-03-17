package main

import (
	"github.com/SongCastle/KoR/api/router"
	"github.com/SongCastle/KoR/internal/encryptor"
	"github.com/SongCastle/KoR/internal/jwt"
	"github.com/SongCastle/KoR/volume/db"
)

func setUp() error {
	if err := db.SetUp(); err != nil {
		return err
	}
	encryptor.Init()
	jwt.Init()
	return nil
}

func main() {
	if err := setUp(); err != nil {
		println("Failed to launch Server.")
		return
	}
	router.Routes().Run(":8080")
}
