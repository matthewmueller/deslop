package example

type User struct {
	name  string
	bytes []byte
}

func (u *User) GetName() string { // want "method GetName should be renamed to Name; avoid Get prefixes"
	return u.name
}

func (u *User) getName() string {
	return u.name
}

func (u *User) GetBytes() []byte { // want "method GetBytes should be renamed to Bytes; avoid Get prefixes"
	return u.bytes
}

func (u *User) Get() string {
	return u.name
}

func (u *User) Name() string {
	return u.name
}
