package redisclient

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"helpers/utils"
	"sync"
	"time"
)

var redispool_init_ctx sync.Once
var redispool_instance *redis.Pool

func GetRedisPool(redissvr []string, redissvrcnt, conntimeout, readtimeout, writetimeout, maxidle, maxactive int) *redis.Pool {

	redispool_init_ctx.Do(func() {

		redispool_instance = &redis.Pool{
			MaxIdle:   maxidle,
			MaxActive: maxactive,
			Dial: func() (redis.Conn, error) {
				for i := 0; i < redissvrcnt*utils.RANDOM_TRY_MULTIPLE; i++ {
					index := utils.GenRandomInt(redissvrcnt)
					v := redissvr[index]
					c, err := redis.DialTimeout("tcp", v, time.Duration(conntimeout)*time.Millisecond, time.Duration(readtimeout)*time.Millisecond, time.Duration(writetimeout)*time.Millisecond)
					if err == nil && c != nil {
						return c, nil
					}

				}

				return nil, errors.New("redispool: cannot connect to any redis server")
			},
		}
	})

	return redispool_instance
}
