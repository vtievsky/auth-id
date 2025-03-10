package tarantoolclient

type RolePrivilege struct {
	RoleID      uint64 `json:"role_id"`
	PrivilegeID uint64 `json:"privilege_id"`
	Allowed     bool   `json:"allowed"`
}

type RolePrivilegeCreated struct {
	RoleID      uint64 `json:"role_id"`
	PrivilegeID uint64 `json:"privilege_id"`
	Allowed     bool   `json:"allowed"`
}

func (s RolePrivilegeCreated) ToTuple() Tuple {
	return Tuple{
		s.RoleID,
		s.PrivilegeID,
		s.Allowed,
	}
}

type RolePrivilegeUpdated struct {
	RoleID      uint64 `json:"role_id"`
	PrivilegeID uint64 `json:"privilege_id"`
	Allowed     bool   `json:"allowed"`
}

func (s RolePrivilegeUpdated) ToTuple() Tuple {
	return Tuple{
		s.RoleID,
		s.PrivilegeID,
		s.Allowed,
	}
}

type RolePrivilegeDeleted struct {
	RoleID      uint64 `json:"role_id"`
	PrivilegeID uint64 `json:"privilege_id"`
}

func (s RolePrivilegeDeleted) ToTuple() Tuple {
	return Tuple{
		s.RoleID,
		s.PrivilegeID,
	}
}
