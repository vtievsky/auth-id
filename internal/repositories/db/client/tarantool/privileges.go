package clienttarantool

type Privilege struct {
	ID          uint64 `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s Tuple) ToPrivilege() Privilege {
	return Privilege{
		ID:          s[0].(uint64), //nolint:forcetypeassert
		Code:        s[1].(string), //nolint:forcetypeassert
		Name:        s[2].(string), //nolint:forcetypeassert
		Description: s[3].(string), //nolint:forcetypeassert
	}
}
