package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/SongCastle/KoR/api/model"
	"github.com/SongCastle/KoR/internal/ecode"
	"github.com/SongCastle/KoR/internal/jwt"
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
		// TODO: 期限切れの JWT Token について
		jToken := strings.TrimSpace(auth[authLen - 1])
		rjt, err := jwt.Verify(jToken)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err).SetMeta(ecode.CodeJson("FailToValidateAuthToken"))
			return
		}
		// JWT Token からユーザを検索
		uParams := &model.UserParams{ID: rjt.UserID}
		user, err := model.GetUser(
			model.SelectColumns("id", "login", "email"),
			model.WhereUser(uParams),
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ecode.CodeJson("InvalidAuthToken"))
			return
		}
		tParams := &model.TokenParams{UserID: rjt.UserID, UUID: &rjt.ID}
		token, err := model.GetToken(model.WhereToken(tParams))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, ecode.CodeJson("InvalidAuthToken"))
			return
		}
		user.CurrentToken = token
		c.Set("CurrentUser", user)
		log.Printf("[DEBUG] User#%d (%s)", user.ID, user.Login)
		c.Header(TokenHeader, jToken)

		c.Next()
	}
}
