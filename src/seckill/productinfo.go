//write product information to redis hash map, and get the information from it
package seckill

import (
	"errors"
	//	"fmt"
	logger "github.com/xlog4go"
	"helpers/iowrapper"
	"strconv"
	"strings"
)

//const (
//	PRODUCTINFO = "product_info"
//)

type ProductInfo struct {
	Pid            string
	Pnum           int64
	Seckillingtime string
}

//merge the values of value field to a string
func mergeProductInfo(pi ProductInfo) (string, error) {
	mergeStr := (strconv.FormatInt(pi.Pnum, 10)) + "|" + (pi.Seckillingtime)
	return mergeStr, nil
}

//unmerge a string to the kinds of values
func unmergeProductInfo(pid string, str string) (ProductInfo, error) {
	tempStrs := strings.Split(str, "|")
	//	fmt.Printf("\ntmpstr:%+v\n", len(tempStrs))
	if len(tempStrs) < 2 {
		//		fmt.Printf("\nstr:%+v\n", str)
		return ProductInfo{}, errors.New("product info error")
	}
	pnum_int64, _ := strconv.ParseInt(tempStrs[0], 10, 64)
	pi := ProductInfo{Pid: pid, Pnum: pnum_int64, Seckillingtime: tempStrs[1]}

	return pi, nil
}

//set product infomation
func SetProductInfo(pi ProductInfo, redisCli *iowrapper.RedisClient) error {
	key := PRODUCTINFO
	var values []interface{}
	values = append(values, pi.Pid)
	tempValue, _ := mergeProductInfo(pi)
	values = append(values, tempValue)

	if redisCli != nil {
		_, err := redisCli.Hmset(key, values)

		if err != nil {
			logger.Error("error=[redis_hmset_hash_failed] pid=[%s] message=[%s] err=[%s]", key, pi.Pid, tempValue, err.Error())
		}
		return err
	}
	return nil
}

//Get product information from a redis hash map
func GetProductInfo(pid string, redisCli *iowrapper.RedisClient) (ProductInfo, error) {
	var pi ProductInfo
	tempStr1, err := redisCli.Hget(PRODUCTINFO, pid)
	if err != nil {
		return pi, err
	}
	tempStr := string(tempStr1)
	pi, err = unmergeProductInfo(pid, tempStr)
	return pi, err
}

//Get all product information from a redis hash map
func GetAllProductInfo(redisCli *iowrapper.RedisClient) ([]ProductInfo, error) {
	tempStr, err := redisCli.Hgetall(PRODUCTINFO)
	//fmt.Printf("tmpstr:%+v,err:%v", tempStr, err)

	products := map[string]string{}
	productInfos := []ProductInfo{}

	for i := 0; i < len(tempStr); i = i + 2 {
		products[tempStr[i]] = tempStr[i+1]
	}

	for k, v := range products {
		pi, err := unmergeProductInfo(k, v)
		//fmt.Printf("[pi]:%v,err:%v", pi, err)
		if err == nil {
			productInfos = append(productInfos, pi)
		}
	}

	return productInfos, err
}
