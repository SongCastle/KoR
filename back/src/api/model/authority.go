package model

import "time"

const (
	// Rights
	CREATE = uint8(1)
	READ   = CREATE << 1
	UPDATE = READ   << 1
	DELETE = UPDATE << 1
	ALL    = CREATE << 7
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

type AuthorityBuilder = func(*Authority)

func WithTokenID(tokenID uint64) AuthorityBuilder {
	return func(a *Authority) {
		a.TokenID = tokenID
	}
}

func WithUsersRight() AuthorityBuilder {
	return func(a *Authority) {
		a.Type = UserType
	}
}

func WithCreateRight(all bool) AuthorityBuilder {
	return func(a *Authority) {
		if !a.CanCreate(all) {
			if all {
				a.addAllCRUDRight(CREATE)
				return
			}
			a.addCRUDRight(CREATE)
		}
	}
}

func WithReadRight(all bool) AuthorityBuilder {
	return func(a *Authority) {
		if !a.CanRead(all) {
			if all {
				a.addAllCRUDRight(READ)
				return
			}
			a.addCRUDRight(READ)
		}
	}
}

func WithUpdateRight(all bool) AuthorityBuilder {
	return func(a *Authority) {
		if !a.CanUpdate(all) {
			if all {
				a.addAllCRUDRight(UPDATE)
				return
			}
			a.addCRUDRight(UPDATE)
		}
	}
}

func WithDeleteRight(all bool) AuthorityBuilder {
	return func(a *Authority) {
		if !a.CanDelete(all) {
			if all {
				a.addAllCRUDRight(DELETE)
				return
			}
			a.addCRUDRight(DELETE)
		}
	}
}

func (a *Authority) Build(adds ...AuthorityBuilder) {
	for _, add := range adds {
		add(a)
	}
}

func (a *Authority) IsUserRight() bool {
	return a.Type == UserType
}

func (a *Authority) CanCreate(all bool) bool {
	return a.can(CREATE) && (!all || a.canAll())
}

func (a *Authority) CanRead(all bool) bool {
	return a.can(READ) && (!all || a.canAll())
}

func (a *Authority) CanUpdate(all bool) bool {
	return a.can(UPDATE) && (!all || a.canAll())
}

func (a *Authority) CanDelete(all bool) bool {
	return a.can(DELETE) && (!all || a.canAll())
}

func (a *Authority) can(v uint8) bool {
	return a.Right != nil && a.Right[7] & v == v
}

func (a *Authority) canAll() bool {
	return a.Right != nil && a.Right[0] & ALL == ALL
}

func (a *Authority) rights() []uint8 {
	if a.Right == nil {
		a.Right = []uint8{0, 0, 0, 0, 0, 0, 0, 0}
	}
	return a.Right
}

func (a *Authority) addCRUDRight(v uint8) {
	r := a.rights()
	r[7] += v
}

func (a *Authority) addAllCRUDRight(v uint8) {
	r := a.rights()
	r[0] += v
}
