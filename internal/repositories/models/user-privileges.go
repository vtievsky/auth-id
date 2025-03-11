// Модели для обмена данными между слоями
package models

import "time"

type UserPrivilege struct {
	UserID      uint64
	PrivilegeID uint64
	DateIn      time.Time
	DateOut     time.Time
}
