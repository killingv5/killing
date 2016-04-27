package seckill

import (
	"testing"
	//"helpers/iowrapper"
	"fmt"
)

/*func TestDealRequestQueue(t *testing.T) {
	redisCli := &iowrapper.RedisClient{
		Servers:        []string{"127.0.0.1:6379"},
	}
	err := redisCli.Init()
	fmt.Println(err)
	DealRequestQueue(111, redisCli)
}*/

func TestConvertBytesToInt(t *testing.T) {
	b := []byte{0x00, 0x00, 0x03, 0xe8}
	fmt.Println(ConvertBytesToInt(b))
}
