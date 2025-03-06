package dbclient

import (
	"context"
	"fmt"
	"strings"
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
			ID:       2, //nolint:mnd
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

	_, err := s.getUser(user.Login)
	if err != nil {
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

	return nil, ErrUserAlreadyExists
}

func (s *Stub) UpdateUser(ctx context.Context, user UserUpdated) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	usr, err := s.getUser(user.Login)
	if err != nil {
		return nil, err
	}

	usr.FullName = user.FullName
	usr.Blocked = user.Blocked

	return usr, nil
}

func (s *Stub) DeleteUser(ctx context.Context, login string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.getUser(login)
	if err != nil {
		return err
	}

	ul := make([]*User, 0, len(s.users)-1)

	for _, k := range s.users {
		if strings.EqualFold(k.Login, login) {
			continue
		}

		ul = append(ul, k)
	}

	s.users = ul

	return nil
}

func (s *Stub) getUser(login string) (*User, error) {
	for _, k := range s.users {
		if strings.EqualFold(k.Login, login) {
			return k, nil
		}
	}

	return nil, ErrUserNotFound
}
