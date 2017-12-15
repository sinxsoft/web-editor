package models

type UserExt struct {
	RememberMe bool
	IP         string
	U          User
}

func GenUserExt() UserExt {
	u := new(UserExt)
	u.RememberMe = false
	return *u
}
