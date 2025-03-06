package dbclient

import (
	"context"
	"fmt"
)

type Stub struct {
}

func NewStub() *Stub {
	return &Stub{}
}

func (s *Stub) GetUsers(ctx context.Context) ([]*User, error) {
	num := 25
	users := make([]*User, 0, num)

	users = append(users,
		&User{
			ID:       1,
			Login:    "pupkin_vi",
			FullName: "Пупкин Василий Иванович",
			Blocked:  false,
		},
		&User{
			ID:       2,
			Login:    "papiroskina_mn",
			FullName: "Папироскина Мария Николаевна",
			Blocked:  false,
		},
	)

	for k := len(users); k < num; k++ {
		users = append(users,
			&User{
				ID:       k,
				Login:    fmt.Sprintf("user%d", k),
				FullName: fmt.Sprintf("Пользователь%d", k),
				Blocked:  (k%7 == 0),
			},
		)
	}

	return users, nil
}
