package seckill

import (
	"testing"
	"fmt"
	"helpers/iowrapper"
)

func TestDealRequestQueue(t *testing.T) {
	redisCli := &iowrapper.RedisClient{
		Servers:        []string{"127.0.0.1:6379"},
	}
	err := redisCli.Init()
	fmt.Println(err)
	DealRequestQueue(111, 100, redisCli)
}
