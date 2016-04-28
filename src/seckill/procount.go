package seckill

import (
	logger "github.com/xlog4go"
	"helpers/iowrapper"
	"strconv"
)

/**
* 根据商品id查询count值
* pid：商品id
* client：Redis实例
* return (count<int64>, error) 异常情况下count返回0
**/
func GetProductCount(pid string, client *iowrapper.RedisClient) (int64, error) {
	res, err := client.Get(COUNT_TYPE + pid)
	if err != nil{
		logger.Error("error=[count查询-》失败] key=[%s] err=[%s]", pid, err.Error())
		return 0, err
	}
	if res == nil{
		logger.Error("error=[count查询-》商品不存在] key=[%s] err=[%s]", pid, err.Error())
		return 0, err	
	}
	countInt, _ := strconv.ParseInt(string(res), 10, 64)
	return countInt, nil
}