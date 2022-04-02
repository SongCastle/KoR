package middleware

import (
	"net/http"

	"github.com/SongCastle/KoR/internal/ecode"
	"github.com/SongCastle/KoR/internal/env"
	"github.com/SongCastle/KoR/internal/log"
	"github.com/gin-gonic/gin"
)

var Certification = "cert"

func init() {
	Certification = env.Get("KOR_CERT")
}

func CertMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cert, err := extractCredential(c.Request.Header.Get(AuthorizationHeader))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ecode.CodeJson(err.Error()))
			return
		}
		if cert != Certification {
			log.Warnf("Invalid cert access: %s", c.ClientIP())
			c.AbortWithStatusJSON(http.StatusUnauthorized, ecode.CodeJson("InvalidCertification"))
			return
		}

		c.Next()
	}
}
