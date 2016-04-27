package seckill

import (
	"helpers/iowrapper"
	logger "github.com/xlog4go"
	"strconv"
	_ "fmt"
)	

type UserSeckingInfo struct {
	Status int64
	Goodsid int64
} 

func queryUserSeckillingInfo(uid string, pid string, client *iowrapper.RedisClient) (*UserSeckingInfo, error) {
	res, err := client.Hget(pid, uid)
	if res == "" {
		logger.Warn("秒杀失败", err.Error())
		return &UserSeckingInfo{Status:2, Goodsid:0}, err
	}
	
	gid, err := strconv.ParseInt(res, 10, 64)
	if gid == 0 {
		logger.Warn("秒杀中", err.Error())
		return &UserSeckingInfo{Status:3, Goodsid:0}, err
	}
	return &UserSeckingInfo{Status:1, Goodsid: gid}, nil
}