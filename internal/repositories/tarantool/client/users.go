// Модели для обмена данными между хранилищем и приложением
package tarantoolclient

type User struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Login   string `json:"login"`
	Blocked bool   `json:"blocked"`
}

type UserCreated struct {
	Name    string `json:"name"`
	Login   string `json:"login"`
	Blocked bool   `json:"blocked"`
}

type UserUpdated struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Login   string `json:"login"`
	Blocked bool   `json:"blocked"`
}
