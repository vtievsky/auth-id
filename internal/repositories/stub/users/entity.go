package stubusers

type User struct {
	ID       int
	Login    string
	FullName string
	Blocked  bool
}

type UserCreated struct {
	Login    string
	FullName string
	Blocked  bool
}

type UserUpdated struct {
	Login    string
	FullName string
	Blocked  bool
}
