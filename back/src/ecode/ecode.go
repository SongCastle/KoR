package ecode

import "github.com/gin-gonic/gin"

var codeTable map[string]string = map[string]string{
	"FailToGetUsers"         : "fail_to_get_users",
	"BlankUserID"            : "blank_user_id",
	"InvalidUserID"          : "invalid_user_id_format",
	"FailToGetUser"          : "fail_to_get_user",
	"InvalidCreateUserParams": "invalid_create_user_params",
	"FailToCreateUser"       : "fail_to_create_user",
	"BlankAuthHeader"        : "blank_authorization_header",
	"InvalidAuthHeader"      : "Invalid_authorization_header_format",
	"FailToValidateAuthToken": "fail_to_validate_authorization_token",
	"InvalidAuthToken"       : "invalid_authorization_token_format",
	"InvalidUpdateUserParams": "invalid_update_user_params",
	"FailToUpdateUser"       : "fail_to_update_user",
	"FailToDeleteUser"       : "fail_to_delete_user",
	"InvalidAuthUserParams"  : "invalid_authorize_user_params",
	"BlankLogin"             : "blank_login",
	"BlankPassword"          : "blank_password",
	"FailToAuth"             : "invalid_login_or_password",
	"FailToGenerateAuthToken": "fail_to_generate_authorization_token",
}

func CodeJson(key string) gin.H {
	code, ok := codeTable[key]
	if !ok {
		code = "unknown"
	}
	return gin.H{"code": code}
}
