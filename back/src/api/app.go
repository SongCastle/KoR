package api

import (
	"net/http"

	"github.com/SongCastle/KoR/ecode"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func NoRoute(c *gin.Context) {
	c.String(http.StatusNotFound, "Invalid Path (%s)", c.Request.URL.Path)
}

func abortWithError(c *gin.Context, status int, code string, err error) {
	c.AbortWithError(status, err).SetMeta(ecode.CodeJson(code))
}

func abortWithJSON(c *gin.Context, status int, code string) {
	c.AbortWithStatusJSON(status, ecode.CodeJson(code))
}
