// Модели для обмена данными между хранилищем и приложением
package tarantoolclient

type Role struct {
	ID          uint64 `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Blocked     bool   `json:"blocked"`
}

type RoleCreated struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Blocked     bool   `json:"blocked"`
}

func (s RoleCreated) ToTuple() Tuple {
	return Tuple{
		nil,
		s.Code,
		s.Name,
		s.Description,
		s.Blocked,
	}
}

type RoleUpdated struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Blocked     bool   `json:"blocked"`
}

func (s RoleUpdated) ToTuple() Tuple {
	return Tuple{
		s.ID,
		s.Name,
		s.Description,
		s.Blocked,
	}
}

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
