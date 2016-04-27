package seckill

import (
	"fmt"
	"helpers/iowrapper"
	"testing"
)

func TestHello(t *testing.T) {
	fmt.Println("hello")
	var client *iowrapper.RedisClient = &iowrapper.RedisClient{Servers:[]string{"127.0.0.1:6379"}}

	client.Init()
	err := Pushtoredis("111", "00001", client)
	//products, err := client.Hgetall("product_111")
	fmt.Println(err)
}
