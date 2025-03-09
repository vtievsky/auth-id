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

type RoleUpdated struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Blocked     bool   `json:"blocked"`
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
