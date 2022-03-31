package v1

import (
	"net/http"
	"strconv"

	"github.com/SongCastle/KoR/api/middleware"
	"github.com/SongCastle/KoR/api/model"
	"github.com/SongCastle/KoR/internal/jwt"
	"github.com/gin-gonic/gin"
)

var responseKeys = []string{"id", "login", "email", "created_at", "updated_at"}

func ShowUsers(c *gin.Context) {
	// 権限確認
	cUser, err := currentUser(c)
	if err != nil {
		abortWithJSON(c, http.StatusBadRequest, "UnidentifiedUser")
		return
	}
	if !cUser.CurrentToken().UserAuthority().CanRead(true) {
		abortWithJSON(c, http.StatusUnauthorized, "NotPermitted")
		return
	}
	// ユーザ取得
	users, err := model.GetUsers(
		model.SelectColumns(responseKeys...),
	)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToGetUsers", err)
		return
	}
	c.JSON(http.StatusOK, &users)
}

func ShowUser(c *gin.Context) {
	sid := c.Param("id")
	if sid == "" {
		abortWithJSON(c, http.StatusBadRequest, "BlankUserID")
		return
	}
	id, err := strconv.ParseUint(sid, 10, 64)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "InvalidUserID", err)
		return
	}
	// 権限確認
	cUser, err := currentUser(c)
	if err != nil {
		abortWithJSON(c, http.StatusBadRequest, "UnidentifiedUser")
		return
	}
	if !cUser.CurrentToken().UserAuthority().CanRead(id != cUser.ID) {
		abortWithJSON(c, http.StatusUnauthorized, "NotPermitted")
		return
	}
	// ユーザ取得
	params := &model.UserParams{ID: id}
	user, err := model.GetUser(
		model.SelectColumns(responseKeys...),
		model.WhereUser(params),
	)
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
	user := model.NewUser(&userParams)
	if err := user.Create(); err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToCreateUser", err)
		return
	}

	if err := user.CreateToken(false); err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToGenerateAuthToken", err)
		return
	}
	// JWT Token 生成
	jt, err := jwt.Generate(user.CurrentToken().UUID, user.Login, user.ID)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToGenerateAuthToken", err)
		return
	}
	// Token をヘッダへセット
	c.Header(middleware.TokenHeader, jt.Token)
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	// 対象の User ID を取得
	sid := c.Param("id")
	if sid == "" {
		abortWithJSON(c, http.StatusBadRequest, "BlankUserID")
		return
	}
	id, err := strconv.ParseUint(sid, 10, 64)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "InvalidUserID", err)
		return
	}
	// 権限確認
	cUser, err := currentUser(c)
	if err != nil {
		abortWithJSON(c, http.StatusBadRequest, "UnidentifiedUser")
		return
	}
	if !cUser.CurrentToken().UserAuthority().CanUpdate(id != cUser.ID) {
		abortWithJSON(c, http.StatusUnauthorized, "NotPermitted")
		return
	}
	// ユーザ更新
	uParams := model.UserParams{ID: id}
	user, err := model.GetUser(model.WhereUser(&uParams))
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToUpdateUser", err)
		return
	}
	if err := c.ShouldBindJSON(&uParams); err != nil {
		abortWithError(c, http.StatusBadRequest, "InvalidUpdateUserParams", err)
		return
	}
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToUpdateUser", err)
		return
	}
	if isBlank(uParams.Login) {
		abortWithJSON(c, http.StatusBadRequest, "BlankLogin")
		return
	}
	if isBlank(uParams.Password) {
		abortWithJSON(c, http.StatusBadRequest, "BlankPassword")
		return
	}
	user.BindParams(&uParams)
	user.ID = id
	if err := user.Update(); err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToUpdateUser", err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	sid := c.Param("id")
	if sid == "" {
		abortWithJSON(c, http.StatusBadRequest, "BlankUserID")
		return
	}
	id, err := strconv.ParseUint(sid, 10, 64)
	if err != nil {
		abortWithError(c, http.StatusBadRequest, "InvalidUserID", err)
		return
	}
	// 権限確認
	cUser, err := currentUser(c)
	if err != nil {
		abortWithJSON(c, http.StatusBadRequest, "UnidentifiedUser")
		return
	}
	if !cUser.CurrentToken().UserAuthority().CanDelete(id != cUser.ID) {
		abortWithJSON(c, http.StatusUnauthorized, "NotPermitted")
		return
	}
	// ユーザ削除
	user := model.NewUser(&model.UserParams{ID: id})
	if err := user.Delete(); err != nil {
		abortWithError(c, http.StatusBadRequest, "FailToDeleteUser", err)
		return
	}
	c.Status(http.StatusNoContent)
}
