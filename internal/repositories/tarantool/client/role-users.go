package tarantoolclient

import "time"

type RoleUser struct {
	RoleID  uint64    `json:"role_id"`
	UserID  uint64    `json:"user_id"`
	DateIn  time.Time `json:"date_in"`
	DateOut time.Time `json:"date_out"`
}

type RoleUserCreated struct {
	RoleID  uint64    `json:"role_id"`
	UserID  uint64    `json:"user_id"`
	DateIn  time.Time `json:"date_in"`
	DateOut time.Time `json:"date_out"`
}

func (s RoleUserCreated) ToTuple() Tuple {
	return Tuple{
		s.RoleID,
		s.UserID,
		s.DateIn.Unix(),
		s.DateOut.Unix(),
	}
}

type RoleUserUpdated struct {
	RoleID  uint64    `json:"role_id"`
	UserID  uint64    `json:"user_id"`
	DateIn  time.Time `json:"date_in"`
	DateOut time.Time `json:"date_out"`
}

func (s RoleUserUpdated) ToTuple() Tuple {
	return Tuple{
		s.RoleID,
		s.UserID,
		s.DateIn.Unix(),
		s.DateOut.Unix(),
	}
}

type RoleUserDeleted struct {
	RoleID uint64 `json:"role_id"`
	UserID uint64 `json:"user_id"`
}

func (s RoleUserDeleted) ToTuple() Tuple {
	return Tuple{
		s.RoleID,
		s.UserID,
	}
}
