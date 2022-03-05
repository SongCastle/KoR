package api

import (
	"net/http"
	"strconv"

	"github.com/SongCastle/KoR/lib/jwt"
	"github.com/SongCastle/KoR/model"
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
	var newUser model.NewUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		abortWithError(c, http.StatusBadRequest, "InvalidCreateUserParams", err)
		return
	}
	user, err := model.CreateUser(&newUser)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToCreateUser", err)
		return
	}
	// TODO: 返却する field の選択
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
	userParams := model.User{}
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
		[]string{"id", "encrypted_password", "login", "password_salt"},
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
	token, err := jwt.Generate(user.ID, user.Login)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToGenerateAuthToken", err)
		return
	}
	c.String(http.StatusOK, token)
}

func responseKeys() []string {
	return []string{"id", "login", "email"}
}
