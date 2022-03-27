package model

import "time"

const (
	CREATE = uint8(1)
	READ   = CREATE << 1
	UPDATE = READ   << 1
	DELETE = UPDATE << 1
	ALL    = CREATE << 7
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

func (a *Authority) BindParams(params *AuthorityParams) {
	bindParamsToObject(params, a)
}

func (a *Authority) can(v uint8) bool {
	return a.Right != nil && a.Right[7] & v == v
}

func (a *Authority) canAll() bool {
	return a.Right != nil && a.Right[0] & ALL == ALL
}
