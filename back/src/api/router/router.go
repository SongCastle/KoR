package router

import (
	h0 "github.com/SongCastle/KoR/api/handler/v0"
	h1 "github.com/SongCastle/KoR/api/handler/v1"
	"github.com/SongCastle/KoR/api/middleware"
	"github.com/gin-gonic/gin"
)

func Routes() *gin.Engine {
	// TODO: 環境変数化など
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	middleware.UseCorsMiddleware(func(corsMiddleware gin.HandlerFunc){
		r.Use(corsMiddleware)
	})

	r.GET("/ping", h0.Ping)

	v1 := r.Group("/v1")
	// Users API
	v1.Use(middleware.ErrorHandleMiddleware())
	{
		usersAPI := v1.Group("/users")

		usersAPI.POST("", h1.CreateUser)
		usersAPI.POST("/token", h1.CreateToken)

		// Require Authorization
		authUsersAPI := usersAPI.Group("")
		authUsersAPI.Use(middleware.AuthHandleMiddleware())
		{
			authUsersAPI.GET("", h1.ShowUsers)
			authUsersAPI.GET("/:id", h1.ShowUser)
			authUsersAPI.PUT("/:id", h1.UpdateUser)
			authUsersAPI.DELETE("/:id", h1.DeleteUser)
			// TODO: 認証トークンの是非を問わず 204 を返却した方が良いかも ...
			authUsersAPI.DELETE("/token", h1.DeleteToken)
		}

		// Admin
		adminUsersAPI := usersAPI.Group("/admin")
		adminUsersAPI.Use(middleware.CertMiddleware())
		{
			adminUsersAPI.POST("/token", h1.CreateAdminToken)
		}
	}

	r.NoRoute(h0.NoRoute)

	return r
}
