// Модели для обмена данными между слоями
package models

type Role struct {
	ID          int
	Code        string
	Name        string
	Description string
	Blocked     bool
}

type RoleCreated struct {
	Name        string
	Description string
	Blocked     bool
}

type RoleUpdated struct {
	Code        string
	Name        string
	Description string
	Blocked     bool
}
