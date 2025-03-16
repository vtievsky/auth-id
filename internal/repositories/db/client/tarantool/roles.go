// Модели для обмена данными между хранилищем и приложением
package clienttarantool

type Role struct {
	ID          uint64 `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Blocked     bool   `json:"blocked"`
}

func (s Tuple) ToRole() Role {
	return Role{
		ID:          s[0].(uint64), //nolint:forcetypeassert
		Code:        s[1].(string), //nolint:forcetypeassert
		Name:        s[2].(string), //nolint:forcetypeassert
		Description: s[3].(string), //nolint:forcetypeassert
		Blocked:     s[4].(bool),   //nolint:forcetypeassert
	}
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
