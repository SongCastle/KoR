package main

import (
	"os"

	"github.com/SongCastle/KoR/api"
	"github.com/SongCastle/KoR/db"
	"github.com/SongCastle/KoR/lib/encryptor"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := load(); err != nil {
		println("Failed to launch Server.")
		return
	}
	serve()
}

func load() error {
	if err := db.InitDB(); err != nil {
		return err
	}
	encryptor.Init()
	return nil
}

func serve() {
	r := gin.Default()

	if frontHost := os.Getenv("FRONT_HOST"); frontHost != "" {
		r.Use(cors.New(cors.Config{
			AllowOrigins: []string{frontHost},
		}))
	}

	v1 := r.Group("/v1")
	v1.GET("/ping", api.Ping)

	// Users API
	v1.GET("/users", api.ShowUsers)
	v1.GET("/users/:id", api.ShowUser)
	v1.PUT("/users/:id", api.UpdateUser)
	v1.POST("/users", api.CreateUser)
	v1.DELETE("/users/:id", api.DeleteUser)

	r.NoRoute(api.NoRoute)

	r.Run(":8080")
}