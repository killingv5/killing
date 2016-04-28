package seckill

import (
	"helpers/iowrapper"
	logger "github.com/xlog4go"
	"strconv"
	"fmt"
)	

/**
* 用户秒杀结果
* status：秒杀状态
* goodsid：秒杀中的商品编号
**/
type UserSeckingInfo struct {
	Status int64
	Goodsid int64
} 

/**
* 查询用户的秒杀结果
* param: uid  用户id，String类型
* param: pid  商品id，String类型
* param: client Redis客户端实例
* result：两个返回值，UserSeckingInfo和error
**/
func QueryUserSeckillingInfo(uid string, pid string, client *iowrapper.RedisClient) (*UserSeckingInfo, error) {
	res, err := client.Hget(PRODUCT_HASH + pid, uid)

	//秒杀结果中没有(pid,uid)，表示未秒中
	if res == "" {
		if err !=nil {
			logger.Warn("errno=[%s] err=[%s]",ERRNO_SECKILING_FAILED, err.Error())
		}
		//fmt.Println(err)
		return &UserSeckingInfo{Status:2, Goodsid:0}, err
	}
	
	gid, err := strconv.ParseInt(res, 10, 64)
	if err !=nil {
		logger.Error("errno=[%s] err=[%s]",ERRNO_PARSE_FAILED, err.Error())
	}
	//fmt.Println(gid)

	if gid == 0 {//有结果，但是商品编号为0，表示正在秒杀中
		if err !=nil {
			logger.Warn("errno=[%s] err=[%s]",ERRNO_SECKILLING, err.Error())
		}
		//fmt.Println(err)
		return &UserSeckingInfo{Status:3, Goodsid:0}, err
	}
	//秒中
	return &UserSeckingInfo{Status:1, Goodsid: gid}, nil
}