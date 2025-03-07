package redisclient

type User struct {
	ID       int    `redis:"id"`
	Login    string `redis:"login"`
	FullName string `redis:"full_name"`
	Blocked  bool   `redis:"blocked"`
}

type UserCreated struct {
	ID       int    `redis:"id"`
	Login    string `redis:"login"`
	FullName string `redis:"full_name"`
	Blocked  bool   `redis:"blocked"`
}

type UserUpdated struct {
	Login    string `redis:"login"`
	FullName string `redis:"full_name"`
	Blocked  bool   `redis:"blocked"`
}
