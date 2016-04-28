package seckill

import (
	"fmt"
	"helpers/iowrapper"
	"testing"
)

func TestHello(t *testing.T) {
	fmt.Println("hello")
	var client *iowrapper.RedisClient = &iowrapper.RedisClient{Servers: []string{"127.0.0.1:6379"}}

	client.Init()
	pi1 := ProductInfo{"111", 100, "2016-04-29 10:00:00"}
	pi2 := ProductInfo{"222", 121, "2016-04-29 11:00:00"}
	pi3 := ProductInfo{"333", 150, "2016-04-29 12:00:00"}

	err1 := SetProductInfo(pi1, client)
	fmt.Println(err1)
	err2 := SetProductInfo(pi2, client)
	fmt.Println(err2)
	err3 := SetProductInfo(pi3, client)
	fmt.Println(err3)

	pi4, err4 := GetProductInfo("222", client)
	fmt.Println(err4)
	fmt.Printf("prod:%+v,err:%+v\n", pi4, err4)
	pi5, err5 := GetAllProductInfo(client)
	fmt.Printf("[All prods]:%+v,err:%+v\n", pi5, err5)
}
