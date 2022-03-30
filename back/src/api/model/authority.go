package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	// Rights
	CREATE = uint8(1)
	READ   = CREATE << 1
	UPDATE = READ   << 1
	DELETE = UPDATE << 1
	// Types
	UserType = "users"
)

type Authority struct {
	ID        uint64    `json:"id,omitempty",gorm:"primaryKey"`
	TokenID   uint64    `json:"token_id,omitempty"`
	Type      string    `json:"type,omitempty"`
	Right     []uint8   `json:"right,omitempty"`
	Token     Token     `json:"-",gorm:"foreignKey:TokenID"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type AuthorityParams struct {
	ID      uint64
	TokenID uint64  `json:"token_id"`
}

func NewAuthority(params *AuthorityParams) *Authority {
	auth := &Authority{}
	auth.BindParams(params)
	return auth
}

func GetAuthority(queries ...queryFunc) (*Authority, error) {
	var auth Authority
	finisher := func(d *gorm.DB) *gorm.DB {
		return d.First(&auth)
	}
	err := executeQueries(append(queries, finisher)...)
	if err == nil {
		return &auth, nil
	}
	return nil, err
}

func GetAuthorities(queries ...queryFunc) ([]Authority, error) {
	var authes []Authority
	finisher := func(d *gorm.DB) *gorm.DB {
		return d.Find(&authes)
	}
	err := executeQueries(append(queries, finisher)...)
	if err == nil {
		return authes, nil
	}
	return nil, err
}

func WhereAuthority(params *AuthorityParams) queryFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(NewAuthority(params))
	}
}

func (a *Authority) Create() error {
	a.ID = 0
	return executeQueries(func(d *gorm.DB) *gorm.DB {
		return d.Create(a)
	})
}

func (a *Authority) Update() error {
	if !a.IsPersisted() {
		return notPersisted
	}
	return executeQueries(func(d *gorm.DB) *gorm.DB {
		return d.Save(a)
	})
}

func (a *Authority) Delete() error {
	if !a.IsPersisted() {
		return notPersisted
	}
	return executeQueries(func(d *gorm.DB) *gorm.DB {
		return d.Delete(a)
	})
}

func (a *Authority) SetTokenID(tokenID uint64) *Authority {
	a.TokenID = tokenID
	return a
}

func (a *Authority) AsUsersRight() *Authority {
	a.Type = UserType
	return a
}

func (a *Authority) AddCreateRight(all bool) *Authority {
	return a.addCRUDRight(CREATE, all)
}

func (a *Authority) AddReadRight(all bool) *Authority {
	return a.addCRUDRight(READ, all)
}

func (a *Authority) AddUpdateRight(all bool) *Authority {
	return a.addCRUDRight(UPDATE, all)
}

func (a *Authority) AddDeleteRight(all bool) *Authority {
	return a.addCRUDRight(DELETE, all)
}

func (a *Authority) CanCreate(all bool) bool {
	return a.canAction(CREATE, all)
}

func (a *Authority) CanRead(all bool) bool {
	return a.canAction(READ, all)
}

func (a *Authority) CanUpdate(all bool) bool {
	return a.canAction(UPDATE, all)
}

func (a *Authority) CanDelete(all bool) bool {
	return a.canAction(DELETE, all)
}

func (a *Authority) IsPersisted() bool {
	return a.ID != 0
}

func (a *Authority) IsUserRight() bool {
	return a.Type == UserType
}

func (a *Authority) BindParams(params *AuthorityParams) {
	bindParamsToObject(params, a)
}

// Hooks
func (a *Authority) BeforeCreate(tx *gorm.DB) error {
	if a.TokenID == 0 {
		return fmt.Errorf("Empty TokenID")
	}
	if a.Type == "" {
		return fmt.Errorf("Empty Type")
	}
	if a.Right == nil {
		return fmt.Errorf("Empty Right")
	}
	if len(a.Right) != 8 {
		return fmt.Errorf("Invalid Right (len: %d)", len(a.Right))
	}
	return nil
}

func (a *Authority) rights() []uint8 {
	if a.Right == nil {
		a.Right = []uint8{0, 0, 0, 0, 0, 0, 0, 0}
	}
	return a.Right
}

func (a *Authority) addCRUDRight(v uint8, all bool) *Authority {
	if !a.canAction(v, all) {
		r := a.rights()
		if all {
			r[0] += v
		} else {
			r[7] += v
		}
	}
	return a
}

func (a *Authority) can(v uint8) bool {
	return a.Right != nil && a.Right[7] & v == v
}

func (a *Authority) canAll(v uint8) bool {
	return a.Right != nil && a.Right[0] & v == v
}

func (a *Authority) canAction(v uint8, all bool) bool {
	if all {
		return a.canAll(v)
	}
	return a.can(v)
}
