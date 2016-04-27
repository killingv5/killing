package seckill

import (
	logger "github.com/xlog4go"
	"helpers/iowrapper"
	"strconv"
)

type ProductSeckingInfo struct {
	Userid  int64
	Goodsid int64
}

func QueryProductSeckingInfo(pid string, client *iowrapper.RedisClient) (error, []ProductSeckingInfo) {
	res, err := client.Hgetall(pid)
	if err != nil {
		logger.Error("error=[商品不存在] key=[%s] err=[%s]", pid, err.Error())
		return err, nil
	}
	productlist := []ProductSeckingInfo{}
	for i := 0; i < len(res); i = i + 2 {
		name := res[i]
		value := res[i+1]
		nameformate, err := strconv.ParseInt(name, 10, 64)
		valueformate, err := strconv.ParseInt(value, 10, 64)
		if err == nil && valueformate > 0 {
			productlist = append(productlist, ProductSeckingInfo{nameformate, valueformate})
		}
	}
	return nil, productlist
}
