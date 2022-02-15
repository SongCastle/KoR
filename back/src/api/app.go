package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func NoRoute(c *gin.Context) {
	c.String(http.StatusNotFound, "Invalid Path (%s)", c.Request.URL.Path)
}

func abortWithError(c *gin.Context, status int, code string, err error) {
	c.AbortWithError(status, err).SetMeta(gin.H{"code": code})
}

func abortWithJSON(c *gin.Context, status int, code string) {
	c.AbortWithStatusJSON(status, gin.H{"code": code})
}
