package middleware

import (
	"github.com/SongCastle/KoR/internal/log"
	"github.com/gin-gonic/gin"
)

func ErrorHandleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if err := c.Errors.Last(); err != nil {
			log.Errorf("%v", err.Error())
			if code, ok := err.Meta.(gin.H); ok {
				c.JSON(c.Writer.Status(), code)
			}
		}
	}
}
