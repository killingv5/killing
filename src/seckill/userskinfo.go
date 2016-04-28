package seckill

import (
	"helpers/iowrapper"
	logger "github.com/xlog4go"
	"strconv"
	"fmt"
)	

type UserSeckingInfo struct {
	Status int64
	Goodsid int64
} 

func QueryUserSeckillingInfo(uid string, pid string, client *iowrapper.RedisClient) (*UserSeckingInfo, error) {
	res, err := client.Hget(PRODUCT_HASH + pid, uid)
	fmt.Println(res)

	if res == "" {
		if err !=nil {
			logger.Warn("秒杀失败", err.Error())			
		}
		fmt.Println(err)
		return &UserSeckingInfo{Status:2, Goodsid:0}, err
	}
	
	gid, err := strconv.ParseInt(res, 10, 64)
	fmt.Println(gid)
	if gid == 0 {
		if err !=nil {
			logger.Warn("秒杀中", err.Error())			
		}
		fmt.Println(err)
		return &UserSeckingInfo{Status:3, Goodsid:0}, err
	}
	fmt.Println(err)
	return &UserSeckingInfo{Status:1, Goodsid: gid}, nil
}