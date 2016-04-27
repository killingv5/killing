package seckill

import (
	"helpers/iowrapper"
	"testing"
)

func Test_redis(t *testing.T) {
	var redisconsole *iowrapper.RedisClient = &iowrapper.RedisClient{Servers: []string{"127.0.0.1:6379"}}
	redisconsole.Init()
	//redisconsole.Set("111", []byte("5"))
	_, _ = queryProductSeckingInfo("112", redisconsole)
}
