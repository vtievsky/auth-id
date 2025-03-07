package stubusers

import (
	"context"
	"fmt"
	"sync"
)

const (
	num = 25
)

type Users struct {
	mu    sync.Mutex
	users map[string]*User
}

func New() *Users {
	users := make(map[string]*User, num)

	users["pupkin_vi"] = &User{
		ID:       0,
		Login:    "pupkin_vi",
		FullName: "Пупкин Василий Иванович",
		Blocked:  false,
	}

	users["papiroskina_mn"] = &User{
		ID:       1,
		Login:    "papiroskina_mn",
		FullName: "Папироскина Мария Николаевна",
		Blocked:  false,
	}

	var login string

	for k := len(users); k < num; k++ {
		login = fmt.Sprintf("user%d", k)

		users[login] = &User{
			ID:       k,
			Login:    login,
			FullName: fmt.Sprintf("Пользователь%d", k),
			Blocked:  (k%7 == 0),
		}
	}

	return &Users{
		mu:    sync.Mutex{},
		users: users,
	}
}

func (s *Users) GetUser(ctx context.Context, login string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if val, ok := s.users[login]; ok {
		return val, nil
	}

	return nil, ErrUserNotFound
}

func (s *Users) GetUsers(ctx context.Context) ([]*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users := make([]*User, 0, len(s.users))

	for _, val := range s.users {
		users = append(users, val)
	}

	return users, nil
}

func (s *Users) CreateUser(ctx context.Context, user UserCreated) (*User, error) {
	const op = "StubUsers.CreateUser"

	_, err := s.GetUser(ctx, user.Login)
	if err != nil {
		s.mu.Lock()
		defer s.mu.Unlock()

		usr := User{
			ID:       len(s.users),
			Login:    user.Login,
			FullName: user.FullName,
			Blocked:  user.Blocked,
		}

		s.users[user.Login] = &usr

		return &usr, nil
	}

	return nil, fmt.Errorf("failed to create user | %s:%w", op, ErrUserAlreadyExists)
}

func (s *Users) UpdateUser(ctx context.Context, user UserUpdated) (*User, error) {
	const op = "StubUsers.UpdateUser"

	usr, err := s.GetUser(ctx, user.Login)
	if err != nil {
		return nil, fmt.Errorf("failed to update user | %s:%w", op, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	usr.FullName = user.FullName
	usr.Blocked = user.Blocked

	return usr, nil
}

func (s *Users) DeleteUser(ctx context.Context, login string) error {
	const op = "StubUsers.DeleteUser"

	_, err := s.GetUser(ctx, login)
	if err != nil {
		return fmt.Errorf("failed to delete user | %s:%w", op, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.users, login)

	return nil
}
