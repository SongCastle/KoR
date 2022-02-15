package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/SongCastle/KoR/ecode"
	"github.com/SongCastle/KoR/lib/jwt"
	"github.com/SongCastle/KoR/model"
	"github.com/gin-gonic/gin"
)

func ShowUsers(c *gin.Context) {
	users, err := model.GetUsers(responseKeys())
	if err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.FailToGetUsers, err)
		return
	}
	c.JSON(http.StatusOK, &users)
}

func ShowUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		abortWithJSON(c, http.StatusBadRequest, ecode.BlankUserID)
		return
	}
	_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.InvalidUserID, err)
		return
	}
	user, err := model.GetUser(&model.UserGetQuery{ID: _id}, responseKeys())
	if err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.FailToGetUser, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var newUser model.NewUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.InvalidCreateUserParams, err)
		return
	}
	user, err := model.CreateUser(&newUser)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.FailToCreateUser, err)
		return
	}
	// TODO: 返却する field の選択
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	// TODO: API 保護用の middleware を追加したい

	// Authorization ヘッダを確認
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		abortWithJSON(c, http.StatusUnauthorized, ecode.BlankAuthHeader)
		return
	}
	auth := strings.Split(authHeader, "Bearer ")
	authLen := len(auth)
	if authLen < 2 {
		abortWithJSON(c, http.StatusUnauthorized, ecode.InvalidAuthHeader)
		return
	}
	// 対象の User ID を取得
	id := c.Param("id")
	if id == "" {
		abortWithJSON(c, http.StatusBadRequest, ecode.BlankUserID)
		return
	}
	_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.InvalidUserID, err)
		return
	}
	// Authorization Token を取得・検証
	token := strings.TrimSpace(auth[authLen - 1])
	ok, err := jwt.Validate(token, _id)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.FailToValidateAuthToken, err)
		return
	}
	if !ok {
		abortWithJSON(c, http.StatusBadRequest, ecode.InvalidAuthToken)
		return
	}
	// ユーザ更新
	var newUser model.NewUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.InvalidUpdateUserParams, err)
		return
	}
	user, err := model.UpdateUser(_id, &newUser)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.FailToUpdateUser, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		abortWithJSON(c, http.StatusBadRequest, ecode.BlankUserID)
		return
	}
	_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.InvalidUserID, err)
		return
	}
	if err := model.DeleteUser(_id); err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.FailToDeleteUser, err)
		return
	}
	c.String(http.StatusOK, "deleted")
}

// TODO: 他の hundler も同様に型定義した方が良いかも ...
type AuthParams struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

func AuthUser(c *gin.Context) {
	var params AuthParams
	if err := c.ShouldBindJSON(&params); err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.InvalidAuthUserParams, err)
		return
	}

	if params.Login == "" {
		abortWithJSON(c, http.StatusBadRequest, ecode.BlankLogin)
		return
	}
	if params.Password == "" {
		abortWithJSON(c, http.StatusBadRequest, ecode.BlankPassword)
		return
	}

	user, err := model.GetUser(
		&model.UserGetQuery{Login: params.Login},
		[]string{"id", "encrypted_password", "login", "password_salt"},
	)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.FailToAuth, err)
		return
	}
	// Password の検証
	if !user.ValidPassword(params.Password) {
		abortWithJSON(c, http.StatusBadRequest, ecode.FailToAuth)
		return
	}
	// JWT Token 生成
	token, err := jwt.Generate(user.ID, user.Login)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, ecode.FailToGenerateAuthToken, err)
		return
	}
	c.String(http.StatusOK, token)
}

func responseKeys() []string {
	return []string{"id", "login", "email"}
}
