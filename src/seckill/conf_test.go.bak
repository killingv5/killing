package seckill

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	conf := SetConfig("../../conf/killing.conf")
	port := conf.GetValue("killing", "port")
	fmt.Println(port) //root

	productID := conf.GetValue("redis","maxIdle")
	fmt.Println(productID)
	/*
	conf.DeleteValue("database", "username")
	username = conf.GetValue("database", "username")
	if len(username) == 0 {
		fmt.Println("username is not exists") //this stdout username is not exists
	}
	conf.SetValue("database", "username", "widuu")
	username = conf.GetValue("database", "username")
	fmt.Println(username) //widuu

	data := conf.ReadList()
	fmt.Println(data)
	*/
}
