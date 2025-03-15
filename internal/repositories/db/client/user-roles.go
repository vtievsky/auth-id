package tarantoolclient

import "time"

type UserRole struct {
	RoleID  uint64    `json:"role_id"`
	UserID  uint64    `json:"user_id"`
	DateIn  time.Time `json:"date_in"`
	DateOut time.Time `json:"date_out"`
}
