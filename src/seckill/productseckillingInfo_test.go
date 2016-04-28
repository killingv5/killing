package seckill

import (
	"fmt"
	logger "github.com/xlog4go"
	"helpers/iowrapper"
	"testing"
)

func Test_redis(t *testing.T) {
	var redisconsole *iowrapper.RedisClient = &iowrapper.RedisClient{Servers: []string{"127.0.0.1:6379"}}
	redisconsole.Init()
	key := "112"
	err, res := QueryProductSeckingInfo(key, redisconsole)
	if err != nil {
		logger.Error("error=[商品不存在] key=[%s] err=[%s]", key, err.Error())
	}
	fmt.Println(len(res))
	for i := 0; i < len(res); i++ {
		fmt.Println(res[i].Userid, res[i].Goodsid)
	}
}
