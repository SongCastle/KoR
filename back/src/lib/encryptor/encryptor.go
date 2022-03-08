package encryptor

import (
	"errors"
	"os"

	"golang.org/x/crypto/bcrypt"
)

const maxPasswordByteSize = 72

var (
	pepper string
	cost int
)

func Init() {
	if _pepper := os.Getenv("PASSWORD_PEPPER"); pepper != "" {
		pepper = _pepper
	}
	cost = bcrypt.DefaultCost
}

func Digest(password string) (string, error) {
	pw := withPepper(password)
	if tooLongPassword(pw) {
		return "", errors.New("Too Long")
	}
	digest, err := bcrypt.GenerateFromPassword([]byte(pw), cost)
	return string(digest), err
}

func Compare(digest, password string) bool {
	pw := withPepper(password)
	if tooLongPassword(pw) {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(digest), []byte((pw))) == nil
}

func withPepper(password string) string {
	if pepper == "" {
		return password
	}
	return password + pepper
}

// TODO: maxPasswordByteSize を超過できるようにしたい
func tooLongPassword(password string) bool {
	return len(password) > maxPasswordByteSize
}
