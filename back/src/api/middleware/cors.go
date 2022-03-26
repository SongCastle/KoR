package middleware

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func UseCorsMiddleware(use func(gin.HandlerFunc)) {
	if frontHost := os.Getenv("FRONT_HOST"); frontHost != "" {
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
