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
	//redisconsole.Set("111", []byte("5"))
	key := "112"
	err, res = queryProductSeckingInfo(key, redisconsole)
	if err != nil {
		logger.Error("error=[商品不存在] key=[%s] err=[%s]", key, err.Error())
	}
	for i := 0; i < len(res); i++ {
		fmt.Println(res[i].Userid, res[i].Goodsid)
	}
}
