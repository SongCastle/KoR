package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/SongCastle/KoR/model"
	"github.com/SongCastle/KoR/lib/jwt"
	"github.com/gin-gonic/gin"
)

func ShowUsers(c *gin.Context) {
	users, err := model.GetUsers(responseKeys())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &users)
}

func ShowUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Blank ID",
		})
		return
	}
	_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invaid ID",
		})
		return
	}
	user, err := model.GetUser(&model.UserGetQuery{ID: _id}, responseKeys())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var newUser model.NewUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user, err := model.CreateUser(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	// TODO: API 保護用の middleware を追加したい

	// Authorization ヘッダを確認
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	auth := strings.Split(authHeader, "Bearer ")
	authLen := len(auth)
	if authLen < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Authorization Header",
		})
		return
	}
	// 対象の User ID を取得
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Blank ID",
		})
		return
	}
	_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invaid ID",
		})
		return
	}
	// Authorization Token を取得・検証
	token := strings.TrimSpace(auth[authLen - 1])
	ok, err := jwt.Validate(token, _id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Authorization Token",
		})
		return
	}
	// ユーザ更新
	var newUser model.NewUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	user, err := model.UpdateUser(_id, &newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Blank ID",
		})
		return
	}
	_id, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invaid ID",
		})
		return
	}
	if err := model.DeleteUser(_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if params.Login == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Blank Login",
		})
		return
	}
	if params.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Blank Password",
		})
		return
	}

	user, err := model.GetUser(
		&model.UserGetQuery{Login: params.Login},
		[]string{"id", "encrypted_password", "login", "password_salt"},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid",
		})
		return
	}
	// Password の検証
	if !user.ValidPassword(params.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid",
		})
		return
	}
	// JWT Token 生成
	token, err := jwt.Generate(user.ID, user.Login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to generate token",
		})
		return
	}
	c.String(http.StatusOK, token)
}

func responseKeys() []string {
	return []string{"id", "login", "email"}
}
