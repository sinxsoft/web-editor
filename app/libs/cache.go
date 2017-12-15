package libs

import "github.com/sinxsoft/web-editor/app/models"

var (
	MAP map[string]models.User = make(map[string]models.User)
)

func SaveToken(token string, um models.User) error {
	MAP[token] = um
	return nil
}

func GetToken(token string) (models.User, bool) {
	v, ok := MAP[token]
	if !ok {
		return *new(models.User), false
	}
	return v, true
}
