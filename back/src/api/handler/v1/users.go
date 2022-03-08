package v1

import (
	"net/http"
	"strconv"

	"github.com/SongCastle/KoR/api/middleware"
	"github.com/SongCastle/KoR/api/model"
	"github.com/SongCastle/KoR/internal/ecode"
	"github.com/SongCastle/KoR/internal/jwt"
	"github.com/gin-gonic/gin"
)

func ShowUsers(c *gin.Context) {
	users, err := model.GetUsers(responseKeys())
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToGetUsers", err)
		return
	}
	c.JSON(http.StatusOK, &users)
}

func ShowUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		abortWithJSON(c, http.StatusBadRequest, "BlankUserID")
		return
	}
	_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "InvalidUserID", err)
		return
	}
	user, err := model.GetUser(&model.UserGetQuery{ID: _id}, responseKeys())
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToGetUser", err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var userParams model.UserParams
	if err := c.ShouldBindJSON(&userParams); err != nil {
		abortWithError(c, http.StatusBadRequest, "InvalidCreateUserParams", err)
		return
	}
	user, err := model.CreateUser(&userParams)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToCreateUser", err)
		return
	}
	// JWT Token 生成
	jt, err := jwt.Generate(user.AuthUUID, user.Login, user.ID)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToGenerateAuthToken", err)
		return
	}
	if _, err:= model.UpdateUser(&model.UserParams{ID: user.ID, AuthUUID: &jt.ID}); err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToGiveAuthToken", err)
		return
	}
	// Token をヘッダへセット
	c.Header(middleware.TokenHeader, jt.Token)
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	// 対象の User ID を取得
	_id := c.Param("id")
	if _id == "" {
		abortWithJSON(c, http.StatusBadRequest, "BlankUserID")
		return
	}
	id, err := strconv.ParseUint(_id, 10, 64)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "InvalidUserID", err)
		return
	}
	// ユーザ更新
	var userParams model.UserParams
	if err := c.ShouldBindJSON(&userParams); err != nil {
		abortWithError(c, http.StatusBadRequest, "InvalidUpdateUserParams", err)
		return
	}
	userParams.ID = id
	user, err := model.UpdateUser(&userParams)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToUpdateUser", err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		abortWithJSON(c, http.StatusBadRequest, "BlankUserID")
		return
	}
	_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "InvalidUserID", err)
		return
	}
	if err := model.DeleteUser(_id); err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToDeleteUser", err)
		return
	}
	c.Status(http.StatusNoContent)
}

// TODO: 他の hundler も同様に型定義した方が良いかも ...
type AuthParams struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func AuthUser(c *gin.Context) {
	var params AuthParams
	if err := c.ShouldBindJSON(&params); err != nil {
		abortWithError(c, http.StatusBadRequest, "InvalidAuthUserParams", err)
		return
	}

	if params.Login == "" {
		abortWithJSON(c, http.StatusBadRequest, "BlankLogin")
		return
	}
	if params.Password == "" {
		abortWithJSON(c, http.StatusBadRequest, "BlankPassword")
		return
	}

	user, err := model.GetUser(
		&model.UserGetQuery{Login: params.Login},
		[]string{"id", "encrypted_password", "login", "password_salt", "auth_uuid"},
	)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToAuth", err)
		return
	}
	// Password の検証
	if !user.ValidPassword(params.Password) {
		abortWithJSON(c, http.StatusBadRequest, "FailToAuth")
		return
	}
	// JWT Token 生成
	jt, err := jwt.Generate(user.AuthUUID, user.Login, user.ID)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToGenerateAuthToken", err)
		return
	}
	if _, err:= model.UpdateUser(&model.UserParams{ID: user.ID, AuthUUID: &jt.ID}); err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToGiveAuthToken", err)
		return
	}
	c.String(http.StatusOK, jt.Token)
}

func UnauthUser(c *gin.Context) {
	_user, ok := c.Get("CurrentUser")
	if !ok {
		abortWithJSON(c, http.StatusBadRequest, "UnidentifiedUser")
		return
	}
	user, ok := _user.(*model.User)
	if !ok {
		abortWithJSON(c, http.StatusBadRequest, "UnidentifiedUser")
		return
	}
	blankUUID := ""
	if _, err:= model.UpdateUser(&model.UserParams{ID: user.ID, AuthUUID: &blankUUID}); err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToDeleteAuthToken", err)
		return
	}
	c.Status(http.StatusNoContent)
}

func responseKeys() []string {
	return []string{"id", "login", "email"}
}

func abortWithError(c *gin.Context, status int, code string, err error) {
	c.AbortWithError(status, err).SetMeta(ecode.CodeJson(code))
}

func abortWithJSON(c *gin.Context, status int, code string) {
	c.AbortWithStatusJSON(status, ecode.CodeJson(code))
}
