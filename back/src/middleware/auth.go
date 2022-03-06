package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/SongCastle/KoR/ecode"
	"github.com/SongCastle/KoR/lib/jwt"
	"github.com/SongCastle/KoR/model"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization" // Request Header
	TokenHeader = "X-Authorization-Token" // Response Header
)

func AuthHandleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization ヘッダを確認
		authHeader := c.Request.Header.Get(AuthorizationHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ecode.CodeJson("BlankAuthHeader"))
			return
		}
		auth := strings.Split(authHeader, "Bearer ")
		authLen := len(auth)
		if authLen < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ecode.CodeJson("InvalidAuthHeader"))
			return
		}
		// Authorization Token を取得・検証
		token := strings.TrimSpace(auth[authLen - 1])
		rjt, err := jwt.Verify(token)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err).SetMeta(ecode.CodeJson("FailToValidateAuthToken"))
			return
		}
		user, err := model.GetUser(
			&model.UserGetQuery{ID: rjt.UserID},
			[]string{"id", "login", "email", "auth_uuid"},
		)
		if err != nil || user.AuthUUID != rjt.ID {
			c.AbortWithStatusJSON(http.StatusBadRequest, ecode.CodeJson("InvalidAuthToken"))
			return
		}
		c.Set("CurrentUser", user)
		log.Printf("[DEBUG] User#%d (%s)", user.ID, user.Login)
		c.Header(TokenHeader, token)

		c.Next()
	}
}
