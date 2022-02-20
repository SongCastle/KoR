package main

import (
	"os"

	"github.com/SongCastle/KoR/api"
	"github.com/SongCastle/KoR/db"
	"github.com/SongCastle/KoR/lib"
	"github.com/SongCastle/KoR/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := setUp(); err != nil {
		println("Failed to launch Server.")
		return
	}
	serve()
}

func setUp() error {
	if err := db.InitDB(); err != nil {
		return err
	}
	lib.SetUp()
	return nil
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
	v1.Use(middleware.ErrorHandleMiddleware())
	{
		v1.GET("/users", api.ShowUsers)
		v1.GET("/users/:id", api.ShowUser)
		// TODO: ユーザ作成後と合わせて、認証もした方が良いかも ...
		// TODO: 認証する場合、 token からユーザを取得する API が必要になりそう
		v1.POST("/users", api.CreateUser)
		v1.PUT("/users/auth", api.AuthUser)

		// Require Authorization
		auth := v1.Group("/")
		auth.Use(middleware.AuthHandleMiddleware())
		{
			auth.PUT("/users/:id", api.UpdateUser)
			auth.DELETE("/users/:id", api.DeleteUser)
			// TODO: 認証トークンの是非を問わず 204 を返却した方が良いかも ...
			auth.DELETE("/users/auth", api.UnauthUser)
		}
	}

	r.NoRoute(api.NoRoute)

	r.Run(":8080")
}
