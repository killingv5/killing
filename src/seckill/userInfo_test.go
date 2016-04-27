package seckill

import (
	"fmt"
	"helpers/iowrapper"
	"testing"
)

func TestHello2(t *testing.T) {
	fmt.Println("hello")
	var client *iowrapper.RedisClient = &iowrapper.RedisClient{Servers:[]string{"127.0.0.1:6379"}}
	
	client.Init()
	res1, err1 := QueryUserSeckillingInfo("3", "product_111", client) // status>0
	res2, err2 := QueryUserSeckillingInfo("1", "product_111", client) // status=0
	res3, err3 := QueryUserSeckillingInfo("13", "product_111", client) //no data
	// products, err := client.Hget("product_111", "3")
	fmt.Println(res1)
	fmt.Println(res2)
	fmt.Println(res3)
	// fmt.Println(products)
	fmt.Println(err1)
	fmt.Println(err2)
	fmt.Println(err3)
}
