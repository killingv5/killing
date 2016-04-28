package seckill

import (
	"fmt"
	"helpers/iowrapper"
	"testing"
	"time"
)

func TestHello22(t *testing.T) {
	fmt.Println("hello")
	var client *iowrapper.RedisClient = &iowrapper.RedisClient{Servers: []string{"127.0.0.1:6379"}}

	client.Init()
	t2, _ := time.Parse(TIMEFORMAT, "2016-04-29 12:00:00")
	t3, _ := time.Parse(TIMEFORMAT, "2017-04-29 23:12:00")
	pi1 := ProductInfo{"111", 100, time.Now()}
	pi2 := ProductInfo{"222", 121, t2}
	pi3 := ProductInfo{"333", 150, t3}

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
