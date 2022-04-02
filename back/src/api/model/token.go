package model

import (
	"errors"
	"time"

	"github.com/SongCastle/KoR/internal/log"
	"github.com/SongCastle/KoR/internal/random"
	"github.com/jinzhu/gorm"
)

const UUIDLen = 32

type Token struct {
	ID          uint64      `json:"id,omitempty",gorm:"primaryKey"`
	UserID      uint64      `json:"user_id,omitempty"`
	UUID        string      `json:"uuid,omitempty"`
	User        User        `json:"-",gorm:"foreignKey:UserID"`
	Authorities []Authority `json:"-",gorm:"foreignKey:TokenID"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
	UpdatedAt   time.Time   `json:"updated_at,omitempty"`
}

type TokenParams struct {
	ID     uint64
	UserID uint64  `json:"user_id"`
	UUID   *string `json:"uuid"`
}

func NewToken(params *TokenParams) *Token {
	token := &Token{}
	token.BindParams(params)
	return token
}

func GetToken(queries ...queryFunc) (*Token, error) {
	var token Token
	finisher := func(d *gorm.DB) *gorm.DB {
		return d.First(&token)
	}
	err := executeQueries(append(queries, finisher)...)
	if err == nil {
		return &token, nil
	}
	return nil, err
}

func GetTokens(queries ...queryFunc) ([]Token, error) {
	var tokens []Token
	finisher := func(d *gorm.DB) *gorm.DB {
		return d.Find(&tokens)
	}
	err := executeQueries(append(queries, finisher)...)
	if err == nil {
		return tokens, nil
	}
	return nil, err
}

func WhereToken(params *TokenParams) queryFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(NewToken(params))
	}
}

func (t *Token) Create() error {
	t.ID = 0
	return executeQueries(func(d *gorm.DB) *gorm.DB {
		return d.Create(t)
	})
}

func (t *Token) Update() error {
	if !t.IsPersisted() {
		return notPersisted
	}
	return executeQueries(func(d *gorm.DB) *gorm.DB {
		return d.Save(t)
	})
}

func (t *Token) Delete() error {
	if !t.IsPersisted() {
		return notPersisted
	}
	return executeQueries(func(d *gorm.DB) *gorm.DB {
		return d.Delete(t)
	})
}

func (t *Token) IsPersisted() bool {
	return t.ID != 0
}

func (t *Token) AddUserAuthority(all bool) error {
	if !t.IsPersisted() {
		return notPersisted
	}
	auth := t.UserAuthority()
	return auth.
		SetTokenID(t.ID).
		AsUsersRight().
		AddCreateRight(all).
		AddReadRight(all).
		AddUpdateRight(all).
		AddDeleteRight(all).
		CreateOrUpdate()
}

func (t *Token) UserAuthority() *Authority {
	for _, auth := range t.authorities() {
		if auth.IsUserRight() {
			return &auth
		}
	}
	return &Authority{}
}

func (t *Token) BindParams(params *TokenParams) {
	bindParamsToObject(params, t)
}

// Hooks
func (t *Token) BeforeCreate(tx *gorm.DB) error {
	if t.UserID == 0 {
		return errors.New("Empty UserID")
	}
	if t.UUID == "" {
		t.UUID = random.Generate(UUIDLen)
	}
	return nil
}

func (t *Token) BeforeDelete(tx *gorm.DB) error {
	authes, err := GetAuthorities(
		SelectColumns("id"),
		WhereAuthority(&AuthorityParams{TokenID: t.ID}),
	)
	if err != nil {
		return err
	}
	for _, auth := range authes {
		if err := auth.Delete(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Token) authorities() []Authority {
	if t.Authorities == nil {
		authes, err := GetAuthorities(
			WhereAuthority(&AuthorityParams{TokenID: t.ID}),
		)
		if err != nil {
			log.Errorf("Get Token's Authorities: %v", err)
			authes = []Authority{}
		}
		t.Authorities = authes
	}
	return t.Authorities
}
