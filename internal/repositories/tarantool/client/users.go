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

func (s UserCreated) ToTuple() Tuple {
	return Tuple{
		nil,
		s.Name,
		s.Login,
		s.Blocked,
	}
}

type UserUpdated struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Login   string `json:"login"`
	Blocked bool   `json:"blocked"`
}

func (s UserUpdated) ToTuple() Tuple {
	return Tuple{
		s.ID,
		s.Name,
		s.Login,
		s.Blocked,
	}
}
