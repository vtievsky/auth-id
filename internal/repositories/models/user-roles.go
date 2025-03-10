// Модели для обмена данными между слоями
package models

import "time"

type UserRole struct {
	RoleID  uint64
	UserID  uint64
	DateIn  time.Time
	DateOut time.Time
}
