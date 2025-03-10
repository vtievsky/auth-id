// Модели для обмена данными между слоями
package models

import "time"

type RoleUser struct {
	RoleID  int
	UserID  int
	DateIn  time.Time
	DateOut time.Time
}

type RoleUserCreated struct {
	RoleID  int
	UserID  int
	DateIn  time.Time
	DateOut time.Time
}

type RoleUserUpdated struct {
	RoleID  int
	UserID  int
	DateIn  time.Time
	DateOut time.Time
}

type RoleUserDeleted struct {
	RoleID int
	UserID int
}
