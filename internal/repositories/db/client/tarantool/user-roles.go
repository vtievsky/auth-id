package clienttarantool

import "time"

type UserRole struct {
	RoleID  uint64    `json:"role_id"`
	UserID  uint64    `json:"user_id"`
	DateIn  time.Time `json:"date_in"`
	DateOut time.Time `json:"date_out"`
}

func (s Tuple) ToUserRole() UserRole {
	return UserRole{
		RoleID:  s[0].(uint64),                      //nolint:forcetypeassert
		UserID:  s[1].(uint64),                      //nolint:forcetypeassert
		DateIn:  time.Unix(int64(s[2].(uint64)), 0), //nolint:forcetypeassert,gosec
		DateOut: time.Unix(int64(s[3].(uint64)), 0), //nolint:forcetypeassert,gosec
	}
}
