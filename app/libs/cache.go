package libs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sinxsoft/web-editor/app/models"
)

var (
//MAP map[string]models.User = make(map[string]models.User)
//mapUserExts map[string]models.UserExt = make(map[string]models.UserExt)
)

// func SaveToken(token string, um models.User) error {
// 	MAP[token] = um
// 	return nil
// }

// func GetToken(token string) (models.User, bool) {
// 	v, ok := MAP[token]
// 	if !ok {
// 		return *new(models.User), false
// 	}
// 	return v, true
// }

//注意是秒
func SaveToken(token string, um models.UserExt, second int) error {

	b, _ := json.Marshal(um)
	client := CreateClient()
	s := string(b[:])
	//status := client.Set(token, s, mSecond*1000)
	//秒*100000000 纳秒
	status := client.Set(token, s, time.Duration(second*1000000000))
	fmt.Println(status)
	//mapUserExts[token] = um
	return nil
}

func GetToken(token string) (models.UserExt, bool) {
	client := CreateClient()
	sc, e := client.Get(token).Result()
	var j *models.UserExt = new(models.UserExt)
	if e != nil {
		fmt.Println(e)
		return *new(models.UserExt), false
	}
	fmt.Println(sc)
	error := json.Unmarshal([]byte(sc), j)

	if error == nil {
		return *j, true
	} else {
		fmt.Println(error)
		return *new(models.UserExt), false
	}

}
