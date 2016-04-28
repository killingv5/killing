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
	fmt.Println(res1)
	fmt.Println(res2)
	fmt.Println(res3)
	fmt.Println(err1)
	fmt.Println(err2)
	fmt.Println(err3)
}

func TestCleanRedis(t *testing.T){
	var client *iowrapper.RedisClient = &iowrapper.RedisClient{Servers:[]string{"127.0.0.1:6379"}}
	
	client.Init()

	
	err :=CleanProduct("222",client)
	fmt.Println(err)
	
}

func TestGetProcount(t *testing.T){
	var client *iowrapper.RedisClient = &iowrapper.RedisClient{Servers:[]string{"127.0.0.1:6379"}}
	client.Init()
	GetProductCount("33333",client)

}

func TestAddProduct(t *testing.T){
	var client *iowrapper.RedisClient = &iowrapper.RedisClient{Servers:[]string{"127.0.0.1:6379"}}
	client.Init()
	err :=AddProduct("4444", "21", "20160309141711",client)
	fmt.Println(err)
	
}


func TestGetProductList(t *testing.T){
	var client *iowrapper.RedisClient = &iowrapper.RedisClient{Servers:[]string{"127.0.0.1:6379"}}
	client.Init()
	json, err :=GetProductList(client)
	fmt.Println(json)
	fmt.Println(err)

}

