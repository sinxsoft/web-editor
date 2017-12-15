package libs

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/astaxie/beego"
	redis "gopkg.in/redis.v4"
)

var (
	isInit   bool = false
	mutex         = &sync.Mutex{}
	db       int
	addr     string
	password string
)

// 创建 redis 客户端
func CreateClient() *redis.Client {

	if !isInit {
		mutex.Lock()
		if !isInit {
			dbb, err := strconv.Atoi(beego.AppConfig.String("redis.db"))
			if err != nil {
				fmt.Println("redis.db没有配置")
				dbb = 0
			}
			db = dbb
			addr = beego.AppConfig.String("redis.url")
			password = beego.AppConfig.String("redis.password")
		}
		isInit = true
		mutex.Unlock()
		fmt.Println(addr + "|" + string(db) + "|" + password)
	}

	op := new(redis.Options)

	op.Addr = addr
	op.DB = db

	if password != "" {
		op.Password = password
	}

	op.PoolSize = 5

	client := redis.NewClient(op)

	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	//json.Marshal()
	//client.Set("","",超时)
	//client.HMSet
	return client
}

// set(key, value)：给数据库中名称为key的string赋予值value
// get(key)：返回数据库中名称为key的string的value
// getset(key, value)：给名称为key的string赋予上一次的value
// mget(key1, key2,…, key N)：返回库中多个string的value
// setnx(key, value)：添加string，名称为key，值为value
// setex(key, time, value)：向库中添加string，设定过期时间time
// mset(key N, value N)：批量设置多个string的值
// msetnx(key N, value N)：如果所有名称为key i的string都不存在
// incr(key)：名称为key的string增1操作
// incrby(key, integer)：名称为key的string增加integer
// decr(key)：名称为key的string减1操作
// decrby(key, integer)：名称为key的string减少integer
// append(key, value)：名称为key的string的值附加value
// substr(key, start, end)：返回名称为key的string的value的子串
