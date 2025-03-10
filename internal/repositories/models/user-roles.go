// Модели для обмена данными между слоями
package models

import "time"

type UserRole struct {
	RoleID  int
	UserID  int
	DateIn  time.Time
	DateOut time.Time
}
