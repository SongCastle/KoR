package encryptor

import (
	"os"

	"github.com/SongCastle/KoR/internal/encryptor/argon2id"
	"github.com/SongCastle/KoR/internal/encryptor/bcrypt"
)

const UseArgon2id = true

var Pepper string = "pepper1234"

func Init() {
	if pepper := os.Getenv("PASSWORD_PEPPER"); pepper != "" {
		Pepper = pepper
	}
}

func Digest(password string) (string, error) {
	if UseArgon2id {
		return argon2id.Digest(withPepper(password))
	}
	return bcrypt.Digest(withPepper(password))
}

func Compare(digest, password string) bool {
	if argon2id.EncryptedByArgon2id(digest) {
		return argon2id.Compare(digest, withPepper(password))
	}
	return bcrypt.Compare(digest, withPepper(password))
}

// MEMO: salt + pepper の長さは 16 bytes 以上にする
func withPepper(password string) string {
	if Pepper == "" {
		return password
	}
	return password + Pepper
}
