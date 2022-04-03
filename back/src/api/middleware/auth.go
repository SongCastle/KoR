package middleware

import (
	"net/http"

	"github.com/SongCastle/KoR/api/model"
	"github.com/SongCastle/KoR/internal/ecode"
	"github.com/SongCastle/KoR/internal/jwt"
	"github.com/SongCastle/KoR/internal/log"
	"github.com/gin-gonic/gin"
)

const TokenHeader = "X-Authorization-Token" // for response header

func AuthHandleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization ヘッダを確認
		jToken, err := extractCredential(c.Request.Header.Get(AuthorizationHeader))
		if err != nil {
			// TODO: status code を内包する error を作成してもよいかも
			c.AbortWithStatusJSON(http.StatusUnauthorized, ecode.CodeJson(err.Error()))
			return
		}
		// Authorization Token を取得・検証
		// TODO: 期限切れの JWT Token について
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
		user.SetCurrentToken(token)
		c.Set("CurrentUser", user)
		log.Debugf("User#%d (%s)", user.ID, user.Login)
		// TODO: Token 削除後のヘッダ
		c.Header(TokenHeader, jToken)

		c.Next()
	}
}
