// Модели для обмена данными между слоями
package models

import "time"

type RoleUser struct {
	RoleID  uint64
	UserID  uint64
	DateIn  time.Time
	DateOut time.Time
}

type RoleUserCreated struct {
	RoleID  uint64
	UserID  uint64
	DateIn  time.Time
	DateOut time.Time
}

type RoleUserUpdated struct {
	RoleID  uint64
	UserID  uint64
	DateIn  time.Time
	DateOut time.Time
}

type RoleUserDeleted struct {
	RoleID uint64
	UserID uint64
}
