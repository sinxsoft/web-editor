package libs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/sinxsoft/web-editor/app/models"
)

var (
	OBJECT_KEY = "OBJ_"
)

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
	return status.Err()
}

func DeleteToken(token string) error {
	client := CreateClient()
	status := client.Del(token)
	return status.Err()
}

func GetToken(token string) (models.UserExt, bool) {
	client := CreateClient()
	sc, e := client.Get(token).Result()
	j := new(models.UserExt)
	if e != nil {
		fmt.Println(e)
		return *new(models.UserExt), false
	}
	fmt.Println(sc)
	error := json.Unmarshal([]byte(sc), j)

	if error == nil {
		//如果是rememberMe为true，则延时7天
		//纳秒
		nsecond := 86400 * time.Second
		if j.RememberMe {
			//则延时7天
			nsecond = 7 * nsecond
		}
		client.Expire(token, nsecond)
		return *j, true
	}
	fmt.Println(error)
	return *new(models.UserExt), false

}

func SaveObject(objectId string, data []byte, delaySecond int) error {
	client := CreateClient()
	status := client.Set(OBJECT_KEY+objectId, data, time.Duration(delaySecond)*time.Second)
	return status.Err()
}

func DelObject(objectId string) error {
	client := CreateClient()
	status := client.Del(OBJECT_KEY+objectId)
	return status.Err()
}

func GetObjectAndDelay(objectId string, delaySecond int) ([]byte, error) {
	client := CreateClient()
	bs, e := client.Get(OBJECT_KEY + objectId).Bytes()
	nsecond := time.Duration(delaySecond) * time.Second
	if e == nil && bs != nil {
		go func() {
			result := client.Expire(OBJECT_KEY+objectId, nsecond)
			if result.Err() != nil {
				fmt.Println("expire data fail:" + OBJECT_KEY + objectId + ",秒:" + strconv.Itoa(delaySecond))
			} else {

				fmt.Println("expire data success:" + OBJECT_KEY + objectId + ",秒:" + strconv.Itoa(delaySecond))
			}
		}()
	}
	return bs, e
}

func GetObjectAndCollect(objectId string) ([]byte, error) {
	client := CreateClient()
	bs, e := client.Get(OBJECT_KEY + objectId).Bytes()
	if e == nil && bs != nil {
		go func() {
			result := client.Del(OBJECT_KEY+objectId)
			if result.Err() != nil {
				fmt.Println("collect data fail:" + OBJECT_KEY + objectId )
			} else {

				fmt.Println("collect data success:" + OBJECT_KEY + objectId )
			}
		}()
	}
	return bs, e
}