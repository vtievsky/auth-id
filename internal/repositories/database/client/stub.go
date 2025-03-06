package dbclient

import (
	"context"
	"fmt"
	"sync"
)

const (
	num = 25
)

type Stub struct {
	mu    sync.Mutex
	users []*User
}

func NewStub() *Stub {
	users := append(make([]*User, 0, num),
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

	return &Stub{
		mu:    sync.Mutex{},
		users: users,
	}
}

func (s *Stub) GetUsers(ctx context.Context) ([]*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.users, nil
}

func (s *Stub) CreateUser(ctx context.Context, user UserCreated) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := len(s.users)
	usr := User{
		ID:       id,
		Login:    user.Login,
		FullName: user.FullName,
		Blocked:  user.Blocked,
	}

	s.users = append(s.users, &usr)

	return &usr, nil
}
