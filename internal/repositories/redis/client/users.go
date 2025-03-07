// Модели для обмена данными между хранилищем и приложением
package redisclient

type User struct {
	ID      int    `redis:"id"`
	Name    string `redis:"name"`
	Login   string `redis:"login"`
	Blocked bool   `redis:"blocked"`
}

type UserCreated struct {
	ID      int    `redis:"id"`
	Name    string `redis:"name"`
	Login   string `redis:"login"`
	Blocked bool   `redis:"blocked"`
}

type UserUpdated struct {
	Name    string `redis:"name"`
	Login   string `redis:"login"`
	Blocked bool   `redis:"blocked"`
}
