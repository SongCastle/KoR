package model

import (
	"errors"
	"time"

	"github.com/SongCastle/KoR/internal/encryptor"
	"github.com/SongCastle/KoR/internal/random"
	"github.com/jinzhu/gorm"
)

// password_salt の長さ
const UserPasswordSaltLen = 32

type User struct {
	ID                uint64    `json:"id,omitempty",gorm:"primaryKey"`
	Login             string    `json:"login,omitempty"`
	Password          string    `json:"-"`
	Email             string    `json:"email,omitempty"`
	PasswordSalt      string    `json:"-,omitempty"`
	EncryptedPassword string    `json:"-,omitempty"`
	CurrentToken      *Token    `json:"-"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

type UserParams struct {
	ID       uint64
	Login    *string `json:"login"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
}

func NewUser(params *UserParams) *User {
	user := &User{}
	user.BindParams(params)
	return user
}

func GetUser(queries ...queryFunc) (*User, error) {
	var user User
	finisher := func(d *gorm.DB) *gorm.DB {
		return d.Omit("Password").First(&user)
	}
	err := executeQueries(append(queries, finisher)...)
	if err == nil {
		return &user, nil
	}
	return nil, err
}

func GetUsers(queries ...queryFunc) ([]User, error) {
	var users []User
	finisher := func(d *gorm.DB) *gorm.DB {
		return d.Omit("Password").Find(&users)
	}
	err := executeQueries(append(queries, finisher)...)
	if err == nil {
		return users, nil
	}
	return nil, err
}

// TODO: WhereXXX を共通化したい / Model(user)
func WhereUser(params *UserParams) queryFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(NewUser(params))
	}
}

func (u *User) Create() error {
	u.ID = 0
	return executeQueries(func(d *gorm.DB) *gorm.DB {
		return d.Omit("Password").Create(u)
	})
}

func (u *User) Update() error {
	if !u.IsPersisted() {
		return notPersisted
	}
	return executeQueries(func(d *gorm.DB) *gorm.DB {
		return d.Omit("Password").Save(u)
	})
}

// TODO: 存在しないユーザの場合について
func (u *User) Delete() error {
	if !u.IsPersisted() {
		return notPersisted
	}
	return executeQueries(func(d *gorm.DB) *gorm.DB {
		return d.Delete(u)
	})
}

func (u *User) CreateToken() error {
	token := NewToken(&TokenParams{UserID: u.ID})
	if err := token.Create(); err != nil {
		return err
	}
	if err := token.AddUserAuthority(false); err != nil {
		return err
	}
	u.CurrentToken = token
	return nil
}

func (u *User) IsPersisted() bool {
	return u.ID != 0
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

func (u *User) ValidateLogin() error {
	if u.Login == "" {
		return errors.New("Blank Login")
	}
	uniq, err := u.isUniqLogin()
	if err != nil {
		return err
	}
	if !uniq {
		return errors.New("Duplicate Login")
	}
	return nil
}

func (u *User) ValidatePassword() error {
	if u.Password == "" {
		return errors.New("Blank Password")
	}
	return nil
}

func (u *User) TestPassword(password string) bool {
	return encryptor.Compare(u.EncryptedPassword, password + u.PasswordSalt)
}

func (u *User) BindParams(params *UserParams) {
	bindParamsToObject(params, u)
}

func (u *User) isUniqLogin() (bool, error) {
	var queries = []queryFunc{
		SelectColumns("id"),
	}
	if u.IsPersisted() {
		queries = append(queries, func(d *gorm.DB) *gorm.DB {
			return d.Where("`login` = ? AND `id` <> ?", u.Login, u.ID).Limit(1)
		})
	} else {

		queries = append(queries, func(d *gorm.DB) *gorm.DB {
			return d.Where("`login` = ?", u.Login).Limit(1)
		})
	}
	users, err := GetUsers(queries...)
	if err == nil {
		return len(users) == 0, nil
	}
	return false, err
}

// Hooks
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if err := u.ValidateLogin(); err != nil {
		return err
	}
	if err := u.ValidatePassword(); err != nil {
		return err
	}
	if err := u.EncryptPassword(); err != nil {
		return err
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if err := u.ValidateLogin(); err != nil {
		return err
	}
	if u.Password != "" {
		if err := u.EncryptPassword(); err != nil {
			return err
		}
	}
	return nil
}

func (u *User) BeforeDelete(tx *gorm.DB) error {
	tokens, err := GetTokens(
		SelectColumns("id"),
		WhereToken(&TokenParams{UserID: u.ID}),
	)
	if err != nil {
		return err
	}
	for _, token := range tokens {
		if err := token.Delete(); err!= nil {
			return err
		}
	}
	return nil
}
