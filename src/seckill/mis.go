package seckill

import (
	logger "github.com/xlog4go"
	"helpers/iowrapper"
)

func CleanProduct(pid string, client *iowrapper.RedisClient) (error){
	res, err :=client.Hkeys(PRODUCT_HASH + pid)
	if err != nil{
		logger.Error("error=[商品清空-》查询失败] key=[%s] err=[%s]", pid, err.Error())
		return err
	}else if res !=nil {
		if len(res) == 0 {
			return nil
		}
		var input_params []interface{}
		for i := 0; i < len(res); i = i + 1 {
			input_params=append(input_params,res[i])
		}

		_, err :=client.Hdel(PRODUCT_HASH + pid, input_params)
		if err != nil {
			logger.Error("error=[商品清空-》清空失败] key=[%s] err=[%s]", pid, err.Error())
			return err
		}
	}
	return nil
}

//func AddProduct(pid string, num) 