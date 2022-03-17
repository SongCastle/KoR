package router

import (
	"os"

	h0 "github.com/SongCastle/KoR/api/handler/v0"
	h1 "github.com/SongCastle/KoR/api/handler/v1"
	"github.com/SongCastle/KoR/api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
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

	r.GET("/ping", h0.Ping)

	v1 := r.Group("/v1")
	// Users API
	v1.Use(middleware.ErrorHandleMiddleware())
	{
		v1.GET("/users", h1.ShowUsers)
		v1.GET("/users/:id", h1.ShowUser)
		v1.POST("/users", h1.CreateUser)
		v1.PUT("/users/auth", h1.AuthUser)

		// Require Authorization
		auth := v1.Group("/")
		auth.Use(middleware.AuthHandleMiddleware())
		{
			auth.PUT("/users/:id", h1.UpdateUser)
			auth.DELETE("/users/:id", h1.DeleteUser)
			// TODO: 認証トークンの是非を問わず 204 を返却した方が良いかも ...
			auth.DELETE("/users/auth", h1.UnauthUser)
		}
	}

	r.NoRoute(h0.NoRoute)

	return r
}
