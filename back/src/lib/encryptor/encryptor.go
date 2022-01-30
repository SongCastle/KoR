package encryptor

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

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
	digest, err := bcrypt.GenerateFromPassword([]byte(withPepper(password)), cost)
	return string(digest), err
}

func Compare(digest, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(digest), []byte((withPepper(password)))) == nil
}

func withPepper(password string) string {
	if pepper == "" {
		return password
	}
	return password + pepper
}
