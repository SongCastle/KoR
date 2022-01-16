package model

import (
	"time"

	"github.com/SongCastle/KoR/db"
)

type NewUser struct {
	Login    string `json:"login,omitempty"`
	// TODO: Encryption
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type User struct {
	NewUser
	ID uint64 `json:"id,omitempty",gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
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
	if err := conn.DB().Create(user).Error; err != nil {
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
	if err := conn.DB().Save(_user).Error; err != nil {
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
