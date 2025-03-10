// Модели для обмена данными между слоями
package models

type RolePrivilege struct {
	RoleID      int
	PrivilegeID int
	Allowed     bool
}

type RolePrivilegeCreated struct {
	RoleID      int
	PrivilegeID int
	Allowed     bool
}

type RolePrivilegeUpdated struct {
	RoleID      int
	PrivilegeID int
	Allowed     bool
}

type RolePrivilegeDeleted struct {
	RoleID      int
	PrivilegeID int
}
