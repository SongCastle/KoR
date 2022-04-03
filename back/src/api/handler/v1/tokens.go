package v1

import (
	"fmt"
	"net/http"

	"github.com/SongCastle/KoR/api/model"
	"github.com/SongCastle/KoR/internal/jwt"
	"github.com/gin-gonic/gin"
)

type tokenParams struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func CreateToken(c *gin.Context) {
	_createToken(false)(c)
}

func CreateAdminToken(c *gin.Context) {
	_createToken(true)(c)
}

func DeleteToken(c *gin.Context) {
	user, err := currentUser(c)
	if err != nil {
		abortWithJSON(c, http.StatusBadRequest, "UnidentifiedUser")
		return
	}
	// TODO: CurrentToken 内のレコード生成有無について
	if err := user.CurrentToken().Delete(); err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToDeleteAuthToken", err)
		return
	}
	c.Status(http.StatusNoContent)
}

func getAuthorizedUser(login, password string) (*model.User, error) {
	if login == "" {
		return nil, fmt.Errorf("BlankLogin")
	}
	if password == "" {
		return nil, fmt.Errorf("BlankPassword")
	}
	params := &model.UserParams{Login: &login}
	user, err := model.GetUser(
		model.SelectColumns("id", "encrypted_password", "login", "password_salt"),
		model.WhereUser(params),
	)
	if err != nil {
		return nil, fmt.Errorf("FailToAuth")
	}
	// Password の検証
	if user.TestPassword(password) {
		return user, nil
	}
	return nil, fmt.Errorf("FailToAuth")
}

func _createToken(admin bool) func(*gin.Context) {
	return func(c *gin.Context) {
		var params tokenParams
		if err := c.ShouldBindJSON(&params); err != nil {
			abortWithError(c, http.StatusBadRequest, "InvalidAuthUserParams", err)
			return
		}
		// ユーザ認証
		user, err := getAuthorizedUser(params.Login, params.Password)
		if err != nil {
			abortWithError(c, http.StatusBadRequest, err.Error(), err)
			return
		}
		if err := user.CreateToken(admin); err != nil {
			abortWithError(c, http.StatusBadRequest, "FailToGenerateAuthToken", err)
			return
		}
		// JWT Token 生成
		jt, err := jwt.Generate(user.CurrentToken().UUID, user.Login, user.ID)
		if err != nil {
			abortWithError(c, http.StatusBadRequest, "FailToGenerateAuthToken", err)
			return
		}
		c.String(http.StatusCreated, jt.Token)
	}
}
