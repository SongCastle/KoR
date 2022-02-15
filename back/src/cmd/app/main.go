package main

import (
	"log"
	"os"

	"github.com/SongCastle/KoR/api"
	"github.com/SongCastle/KoR/db"
	"github.com/SongCastle/KoR/lib/encryptor"
	"github.com/SongCastle/KoR/lib/jwt"
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
	jwt.Init()
	return nil
}

func errorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if err := c.Errors.Last(); err != nil {
			log.Printf("Error: %v\n", err.Error())
			if code, ok := err.Meta.(gin.H); ok {
				c.JSON(c.Writer.Status(), code)
			}
		}
	}
}

func serve() {
	// TODO: 環境変数化など
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	if frontHost := os.Getenv("FRONT_HOST"); frontHost != "" {
		r.Use(cors.New(cors.Config{
			AllowOrigins: []string{frontHost},
			AllowHeaders: []string{
				"Access-Control-Allow-Headers",
				"Content-Type",
			},
			AllowMethods: []string{"PUT", "DELETE"},
		}))
	}

	v1 := r.Group("/v1")
	v1.GET("/ping", api.Ping)

	// Users API
	v1.Use(errorMiddleware())
	{
		v1.GET("/users", api.ShowUsers)
		v1.GET("/users/:id", api.ShowUser)
		v1.PUT("/users/:id", api.UpdateUser)
		v1.POST("/users", api.CreateUser)
		v1.DELETE("/users/:id", api.DeleteUser)
		v1.PUT("/users/auth", api.AuthUser)
	}

	r.NoRoute(api.NoRoute)

	r.Run(":8080")
}
