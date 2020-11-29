package redis

import (
	"log"
	"time"

	redis "github.com/gomodule/redigo/redis"
)

var (
	pool          *redis.Pool
	redisHost     = "127.0.0.1:6379"
	redisPassword = "123456"
)

//newRedisPool 创建redis连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     500,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			//1.打开链接
			conn, err := redis.Dial("tcp", redisHost)
			if err != nil {
				panic(err)
			}

			//2.密码认证
			if _, err := conn.Do("AUTH", redisPassword); err != nil {
				panic(err)
			}
			return conn, nil
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newRedisPool()
	data, err := pool.Get().Do("KEYS", "*")
	if err != nil {
		panic(err)
	}
	log.Printf("Info: redis pool init, all keys:%v", data)
}

func RedisPool() *redis.Pool {
	return pool
}
