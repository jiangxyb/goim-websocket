package Initialize

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"im/global"
	"log"
	"sync"
	"time"
)

var once sync.Once

func InitRedis() {
	addr := fmt.Sprintf("%v:%v", global.SrvConfig.RedisInfo.IP, global.SrvConfig.RedisInfo.Port)
	once.Do(func() {
		global.Redis = redis.NewClient(&redis.Options{
			Network:      "tcp",
			Addr:         addr,
			Password:     global.SrvConfig.RedisInfo.PassWord,
			DB:           0,
			PoolSize:     15,
			MinIdleConns: 10,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			PoolTimeout:  4 * time.Second,

			IdleCheckFrequency: 60 * time.Second,
			IdleTimeout:        5 * time.Minute,
			MaxConnAge:         0 * time.Second,

			MaxRetries:      0,
			MinRetryBackoff: 8 * time.Millisecond,
			MaxRetryBackoff: 512 * time.Millisecond,
		})
		pong, err := global.Redis.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(pong)
	})
}
