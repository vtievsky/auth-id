package clienttarantool

type RolePrivilege struct {
	RoleID      uint64 `json:"role_id"`
	PrivilegeID uint64 `json:"privilege_id"`
	Allowed     bool   `json:"allowed"`
}

func (s Tuple) ToRolePrivilege() RolePrivilege {
	return RolePrivilege{
		RoleID:      s[0].(uint64), //nolint:forcetypeassert
		PrivilegeID: s[1].(uint64), //nolint:forcetypeassert
		Allowed:     s[2].(bool),   //nolint:forcetypeassert
	}
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
