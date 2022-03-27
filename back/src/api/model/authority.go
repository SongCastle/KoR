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
		if all {
			if !a.CanAllCreate() {
				a.addAllCRUDRight(CREATE)
			}
			return
		}
		if !a.CanCreate() {
			a.addCRUDRight(CREATE)
		}
	}
}

func WithReadRight(all bool) AuthorityBuilder {
	return func(a *Authority) {
		if all {
			if !a.CanAllRead() {
				a.addAllCRUDRight(READ)
			}
			return
		}
		if !a.CanRead() {
			a.addCRUDRight(READ)
		}
	}
}

func WithUpdateRight(all bool) AuthorityBuilder {
	return func(a *Authority) {
		if all {
			if !a.CanAllUpdate() {
				a.addAllCRUDRight(UPDATE)
			}
			return
		}
		if !a.CanUpdate() {
			a.addCRUDRight(UPDATE)
		}
	}
}

func WithDeleteRight(all bool) AuthorityBuilder {
	return func(a *Authority) {
		if all {
			if !a.CanAllDelete() {
				a.addAllCRUDRight(DELETE)
			}
			return
		}
		if !a.CanDelete() {
			a.addCRUDRight(DELETE)
		}
	}
}

func (a *Authority) Build(adds ...AuthorityBuilder) {
	for _, add := range adds {
		add(a)
	}
}

func (a *Authority) CanCreate() bool {
	return a.can(CREATE)
}

func (a *Authority) CanRead() bool {
	return a.can(READ)
}

func (a *Authority) CanUpdate() bool {
	return a.can(UPDATE)
}

func (a *Authority) CanDelete() bool {
	return a.can(DELETE)
}

func (a *Authority) CanAllCreate() bool {
	return a.CanCreate() && a.canAll()
}

func (a *Authority) CanAllRead() bool {
	return a.can(READ) && a.canAll()
}

func (a *Authority) CanAllUpdate() bool {
	return a.can(UPDATE) && a.canAll()
}

func (a *Authority) CanAllDelete() bool {
	return a.can(DELETE) && a.canAll()
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
