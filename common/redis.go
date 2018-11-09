package common

import (
	"github.com/go-redis/redis"
	"time"
)

var redisClient *redis.Client
func init()  {
	/*client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6369",
		Password:"",
		Network: "tcp4",
	})
	redisClient = client;
	pong, err := client.Ping().Result()

	if err != nil{
		GetLog().Panicln(pong,"redis 连接异常",err)
	}*/
}


func GetRedisClient()(*redis.Client){

	return redisClient
}

func RedisCmd(){
	// key  value  有效期
	redisClient.Set("ss","ss",1 * time.Second)
	GetLog().Println(redisClient.Get("ss").String())
}
