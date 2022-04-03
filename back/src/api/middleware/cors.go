package middleware

import (
	"github.com/SongCastle/KoR/internal/env"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func UseCorsMiddleware(use func(gin.HandlerFunc)) {
	if frontHost := env.Get("FRONT_HOST"); frontHost != "" {
		use(cors.New(cors.Config{
			AllowOrigins: []string{frontHost},
			AllowHeaders: []string{
				"Access-Control-Allow-Headers",
				"Content-Type",
			},
			AllowMethods: []string{"PUT", "DELETE"},
		}));
	}
}
