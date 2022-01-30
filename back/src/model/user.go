package model

import (
	"errors"
	"math/rand"
	"time"
	"strconv"

	"github.com/SongCastle/KoR/db"
	"github.com/SongCastle/KoR/lib/encryptor"
	"github.com/jinzhu/gorm"
)

// password_salt の長さ
const UserPasswordSaltLen = 32

// Hooks
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Login == "" {
		return errors.New("Blank Login")
	}
	if u.Password == "" {
		return errors.New("Blank Password")
	}
	u.SetPasswordSolt()
	u.EncryptPassword()
	return nil
}

type NewUser struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type User struct {
	NewUser
	ID                uint64 `json:"id,omitempty",gorm:"primaryKey"`
	PasswordSalt      string `json:"password_salt,omitempty"`
	EncryptedPassword string `json:"encrypted_password,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

func (u *User) SetPasswordSolt() {
	u.PasswordSalt = generatePasswordSalt()
}

func (u *User) EncryptPassword() error {
	if u.Password != "" {
		digest, err := encryptor.Digest(u.Password + u.PasswordSalt)
		if err != nil {
			return err
		}
		u.EncryptedPassword = digest
	}
	return nil
}

func (u *User) ValidPassword(password string) bool {
	return encryptor.Compare(u.EncryptedPassword, password + u.PasswordSalt)
}

// TODO: User オブジェクト内に関数群を含めた方がよいかも ...

func GetUsers(cols string) ([]User, error) {
	var users []User

	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := conn.DB().Select(cols).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUser(id uint64, cols string) (*User, error) {
	user := &User{ID: id}

	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return nil, err
	}
	defer conn.Close()

	if err := conn.DB().Select(cols).Take(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(newUser *NewUser) (*User, error) {
	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return nil, err
	}
	defer conn.Close()

	// TODO: 検証
	user := &User{NewUser: *newUser}
	if err := conn.DB().Omit("Password").Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(id uint64, user *NewUser) (*User, error) {
	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return nil, err
	}
	defer conn.Close()

	// TODO: 検証
	_user := &User{ID: id, NewUser: *user}
	if err := _user.EncryptPassword(); err != nil {
		return nil, err
	}

	if err := conn.DB().Model(_user).Omit("Password").Update(_user).Error; err != nil {
		return nil, err
	}
	return _user, nil
}

func DeleteUser(id uint64) error {
	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return err
	}
	defer conn.Close()

	// TODO: 検証
	if err := conn.DB().Delete(&User{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

func generatePasswordSalt() string {
	rand.Seed(time.Now().UnixNano())

	salt := ""
	for i := 0; i < UserPasswordSaltLen; i++ {
		n := rand.Intn(62)
		if n < 26 {
			// 小文字
			salt += string('a' + n)
		} else if n < 52 {
			// 大文字
			salt += string('A' + (n - 26))
		} else {
			// 数字
			salt += strconv.Itoa(n - 52)
		}
	}
	return salt
}
