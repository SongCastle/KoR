package model

import (
	"time"

	"github.com/SongCastle/KoR/lib/encryptor"
	"github.com/SongCastle/KoR/lib/random"
)

// password_salt の長さ
const UserPasswordSaltLen = 32

type User struct {
	ID                uint64 `json:"id,omitempty",gorm:"primaryKey"`
	Login             string `json:"login,omitempty"`
	Password          string `json:"password,omitempty"`
	Email             string `json:"email,omitempty"`
	PasswordSalt      string `json:"-,omitempty"`
	EncryptedPassword string `json:"-,omitempty"`
	AuthUUID          string `json:"-,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

func (u *User) EncryptPassword() error {
	if u.Password == "" {
		return nil
	}
	if u.PasswordSalt == "" {
		// Set PasswordSalt
		u.PasswordSalt = random.Generate(UserPasswordSaltLen)
	}
	// Encrypt Password
	digest, err := encryptor.Digest(u.Password + u.PasswordSalt)
	if err != nil {
		return err
	}
	u.EncryptedPassword = digest
	u.Password = ""
	return nil
}

func (u *User) ValidPassword(password string) bool {
	return encryptor.Compare(u.EncryptedPassword, password + u.PasswordSalt)
}
