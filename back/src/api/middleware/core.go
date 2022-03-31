package middleware

import (
	"errors"
	"strings"
)

const AuthorizationHeader = "Authorization" // for request header

func extractCredential(authHeader string) (string, error) {
	if authHeader == "" {
		return authHeader, errors.New("BlankAuthHeader")
	}
	auth := strings.Split(authHeader, "Bearer ")
	if authLen := len(auth); authLen > 1 {
		return strings.TrimSpace(auth[authLen - 1]), nil
	}
	return "", errors.New("InvalidAuthHeader")
}
