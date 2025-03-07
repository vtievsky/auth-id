// Модели для обмена данными между слоями
package models

type User struct {
	ID      int
	Name    string
	Login   string
	Blocked bool
}

type UserCreated struct {
	Name    string
	Login   string
	Blocked bool
}

type UserUpdated struct {
	Name    string
	Login   string
	Blocked bool
}
