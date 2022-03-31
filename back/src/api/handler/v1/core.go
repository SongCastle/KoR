package v1

import (
	"fmt"

	"github.com/SongCastle/KoR/api/model"
	"github.com/SongCastle/KoR/internal/ecode"
	"github.com/gin-gonic/gin"
)

func abortWithError(c *gin.Context, status int, code string, err error) {
	c.AbortWithError(status, err).SetMeta(ecode.CodeJson(code))
}

func abortWithJSON(c *gin.Context, status int, code string) {
	c.AbortWithStatusJSON(status, ecode.CodeJson(code))
}

func currentUser(c *gin.Context) (*model.User, error) {
	_user, ok := c.Get("CurrentUser")
	if !ok {
		return nil, fmt.Errorf("UnidentifiedUser")
	}
	user, ok := _user.(*model.User)
	if !ok {
		return nil, fmt.Errorf("UnidentifiedUser")
	}
	return user, nil
}

func isBlank(str *string) bool {
	return str != nil && *str == ""
}
