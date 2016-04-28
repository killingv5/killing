package seckill

import (
	logger "github.com/xlog4go"
	"helpers/iowrapper"
	"time"
	"errors"
	// "fmt"
	"encoding/json"
	"strconv"
)

/**
*清空商品的秒杀结果
* param：pid 商品ID
**/
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

/**
* 向Redis添加商品
* pid 商品ID
* num 可供秒杀的商品数量
* timestr 开始秒杀的时间串，格式：yyyyMMddHHmmss
**/
func AddProduct(pid string, num string, timestr string, client *iowrapper.RedisClient) (error){
	//1.add product
	if !checkTime(timestr) {
		return errors.New("error=[添加商品-》日期格式错误]")
	}
	t, _ :=time.ParseInLocation("20060102150405", timestr, time.Local)
	//formatTime := t.Format("2006-01-02 15:04:05")
	numInt, _ :=strconv.ParseInt(num, 10, 64)
	pi := ProductInfo{Pid:pid, Pnum:numInt, Seckillingtime:t}
	SetProductInfo(pi, client)
	//2.add counter
	err := client.Set(COUNT_TYPE + pid, []byte(num))
	if err !=nil {
		logger.Error("error=[添加商品-》添加计算器失败] key=[%s] err=[%s]", pid, err.Error())
		return err
	}
	return nil
}

/**
* 查询商品列表
* result json格式的商品信息
**/
func GetProductList(client *iowrapper.RedisClient) (string, error){
	list, err :=GetAllProductInfo(client)
	if err !=nil {
		return "", err
	}
	b, err := json.Marshal(list)
	if err !=nil{
		logger.Error("error=[查询商品-》查询商品失败] err=[%s]", err.Error())
		return "", err
	}
	return string(b), nil
}


func checkTime(timeStr string) (bool){
	if timeStr == "" || len(timeStr) != 14 {
		return false
	}
	return true
}