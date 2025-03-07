// Модели для обмена данными между слоями
package models

type Role struct {
	ID          int    // Ключевое поле для БД
	Code        string // Ключевое поле для интерфейса
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
	Code        string // Ключевое поле для интерфейса
	Name        string
	Description string
	Blocked     bool
}
