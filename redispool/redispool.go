package redispool

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	Redisclient *redis.Pool
	RedisHost   = "94.191.18.219"
	RedisPort   = "6300"
	RedisDb     = 0
	RedisAuth   = "135a246b"
	NetWork     = "tcp"
)

func init() {
	Redisclient = &redis.Pool{
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial(NetWork, fmt.Sprintf("%s:%s", RedisHost, RedisPort))
			if err != nil {
				return
			}
			_, err = conn.Do("auth", RedisAuth)
			if err != nil {
				conn.Close()
				return
			}
			return conn, nil
		},
		MaxIdle:     16,                //初始连接数
		MaxActive:   10000,             //最大连接数
		IdleTimeout: 300 * time.Second, //超时时间
	}
}
