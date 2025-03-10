// Модели для обмена данными между слоями
package models

type RolePrivilege struct {
	RoleID      uint64
	PrivilegeID uint64
	Allowed     bool
}

type RolePrivilegeCreated struct {
	RoleID      uint64
	PrivilegeID uint64
	Allowed     bool
}

type RolePrivilegeUpdated struct {
	RoleID      uint64
	PrivilegeID uint64
	Allowed     bool
}

type RolePrivilegeDeleted struct {
	RoleID      uint64
	PrivilegeID uint64
}
