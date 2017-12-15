package models

type UserExt struct {
	IP string
	U  User
}

func GenUserExt() UserExt {
	u := new(UserExt)
	return *u
}
