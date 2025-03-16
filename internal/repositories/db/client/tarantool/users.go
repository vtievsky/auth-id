// Модели для обмена данными между хранилищем и приложением
package clienttarantool

type User struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Blocked  bool   `json:"blocked"`
}

func (s Tuple) ToUser() User {
	return User{
		ID:       s[0].(uint64), //nolint:forcetypeassert
		Name:     s[1].(string), //nolint:forcetypeassert
		Login:    s[2].(string), //nolint:forcetypeassert
		Password: s[3].(string), //nolint:forcetypeassert
		Blocked:  s[4].(bool),   //nolint:forcetypeassert
	}
}

type UserCreated struct {
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Blocked  bool   `json:"blocked"`
}

func (s UserCreated) ToTuple() Tuple {
	return Tuple{
		nil,
		s.Name,
		s.Login,
		s.Password,
		s.Blocked,
	}
}

type UserUpdated struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Blocked  bool   `json:"blocked"`
}

func (s UserUpdated) ToTuple() Tuple {
	return Tuple{
		s.ID,
		s.Name,
		s.Login,
		s.Password,
		s.Blocked,
	}
}
