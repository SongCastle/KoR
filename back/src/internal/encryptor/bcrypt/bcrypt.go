package bcrypt

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const maxPasswordByteSize = 72

var (
	Cost         int   = bcrypt.DefaultCost
	TooLongError error = errors.New("Too Long")
)

func Digest(password string) (string, error) {
	if tooLongPassword(password) {
		return "", TooLongError
	}
	digest, err := bcrypt.GenerateFromPassword([]byte(password), Cost)
	if err == nil {
		return string(digest), nil
	}
	return "", err
}

func Compare(digest, password string) bool {
	if tooLongPassword(password) {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(digest), []byte((password))) == nil
}

func tooLongPassword(password string) bool {
	return len(password) > maxPasswordByteSize
}
