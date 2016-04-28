package seckill

import (
	logger "github.com/xlog4go"
	"helpers/iowrapper"
	"strconv"
)

type ProductSeckingInfo struct {
	Userid  int64			`json:"userid"`
	Goodsid int64			`json:"goodsid"`
}

func QueryProductSeckingInfo(pid string, client *iowrapper.RedisClient) (error, []ProductSeckingInfo) {
	res, err := client.Hgetall(PRODUCT_HASH + pid)
	if err != nil {
		logger.Error("errno=[%s] key=[%s] err=[%s]", ERRNO_PRODUCT_NOT_EXIST,pid, err.Error())
		return err, nil
	}
	productlist := []ProductSeckingInfo{}
	for i := 0; i < len(res); i = i + 2 {
		name := res[i]
		value := res[i+1]
		nameformate, err := strconv.ParseInt(name, 10, 64)
		valueformate, err := strconv.ParseInt(value, 10, 64)
		productlist = append(productlist, ProductSeckingInfo{nameformate, valueformate})
		if err != nil{
			logger.Error("parse res failed,%s",err.Error())
		}
	}

	return nil, productlist
}
