package model

import (
	"errors"
	"time"

	"github.com/SongCastle/KoR/db"
	"github.com/SongCastle/KoR/lib/encryptor"
	"github.com/SongCastle/KoR/lib/random"
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

type UserGetQuery struct {
	ID    uint64
	Login string
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
	return nil
}

func (u *User) ValidPassword(password string) bool {
	return encryptor.Compare(u.EncryptedPassword, password + u.PasswordSalt)
}

// TODO: User オブジェクト内に関数群を含めた方がよいかも ...

func GetUsers(cols []string) ([]User, error) {
	var users []User

	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return nil, err
	}
	defer conn.Close()

	cdb := conn.DB()
	if len(cols) > 0 {
		cdb = cdb.Select(cols)
	}
	if err := cdb.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUser(query *UserGetQuery, cols []string) (*User, error) {
	user := User{}
	user.ID, user.Login = query.ID, query.Login

	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return nil, err
	}
	defer conn.Close()

	cdb := conn.DB().Where(&user)
	if len(cols) > 0 {
		cdb = cdb.Select(cols)
	}
	if err := cdb.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
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

func UpdateUser(userParams *User) (*User, error) {
	// TODO: 検証する (引数含む)
	conn := db.NewDB()
	if err := conn.Open(); err != nil {
		return nil, err
	}
	defer conn.Close()
	// TODO: DB Connection を Context で管理したい
	user, err := GetUser(&UserGetQuery{ID: userParams.ID}, []string{"password_salt"})
	if err != nil {
		return nil, err
	}
	// TODO: ゼロ値の識別
	user.Login, user.Email, user.Password = userParams.Login, userParams.Email, userParams.Password
	if err := user.EncryptPassword(); err != nil {
		return nil, err
	}
	if err := conn.DB().Model(&User{ID: userParams.ID}).Omit("Password").Update(user).Error; err != nil {
		return nil, err
	}
	return user, nil
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
