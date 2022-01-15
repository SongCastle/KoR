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
