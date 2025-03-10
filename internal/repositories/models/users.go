// Модели для обмена данными между слоями
package models

type User struct {
	ID      uint64 // Ключевое поле для БД
	Name    string
	Login   string // Ключевое поле для интерфейса
	Blocked bool
}

type UserCreated struct {
	Name    string
	Login   string
	Blocked bool
}

type UserUpdated struct {
	Name    string
	Login   string // Ключевое поле для интерфейса
	Blocked bool
}
