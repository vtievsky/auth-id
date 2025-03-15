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

func (s RoleCreated) ToTuple() Tuple {
	return Tuple{
		nil,
		s.Code,
		s.Name,
		s.Description,
		s.Blocked,
	}
}

type RoleUpdated struct {
	ID          uint64 `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Blocked     bool   `json:"blocked"`
}

func (s RoleUpdated) ToTuple() Tuple {
	return Tuple{
		s.ID,
		s.Code,
		s.Name,
		s.Description,
		s.Blocked,
	}
}
