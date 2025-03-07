// Модели для обмена данными между хранилищем и приложением
package redisclient

type Role struct {
	ID          int    `redis:"id"`
	Code        string `redis:"code"`
	Name        string `redis:"name"`
	Description string `redis:"description"`
	Blocked     bool   `redis:"blocked"`
}

type RoleCreated struct {
	ID          int    `redis:"id"`
	Code        string `redis:"code"`
	Name        string `redis:"name"`
	Description string `redis:"description"`
	Blocked     bool   `redis:"blocked"`
}

type RoleUpdated struct {
	Name        string `redis:"name"`
	Description string `redis:"description"`
	Blocked     bool   `redis:"blocked"`
}
