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
	return &Stub{
		mu:    sync.Mutex{},
		users: make([]*User, 0, num),
	}
}

func (s *Stub) GetUsers(ctx context.Context) ([]*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users = append(s.users,
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

	for k := len(s.users); k < num; k++ {
		s.users = append(s.users,
			&User{
				ID:       k,
				Login:    fmt.Sprintf("user%d", k),
				FullName: fmt.Sprintf("Пользователь%d", k),
				Blocked:  (k%7 == 0),
			},
		)
	}

	return s.users, nil
}
