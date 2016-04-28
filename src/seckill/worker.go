package seckill

import (
	"strconv"
	"helpers/iowrapper"
	logger "github.com/xlog4go"
	"fmt"
)

var PidFlag map[int64]bool

func init() {
	PidFlag = make(map[int64]bool)
}

/**
 *	1.从请求队列中阻塞读请求
 *  2.判断商品当前剩余数量,count==0, 秒杀结束,否则继续
 *  3.userId集合判重
 *  4.写订单信息表
 *  5.商品数量减一
 */
func DealRequestQueue(productId int64, productTotal int64, redisCli *iowrapper.RedisClient)  {
	productType := strconv.FormatInt(productId, 10)
	countType := COUNT_TYPE + productType
	productQueueName := PRODUCT_QUEUE + productType
	userIdSetName := USERID_SET + productType
	productName := PRODUCT_HASH + productType

	for {
		userId, err := redisCli.BLpop(productQueueName, 0)
		if err != nil {
			logger.Error("Error occure in reading queue: %s", err.Error())
		}
		count, _ := redisCli.Get(countType);
		countInt, _ := strconv.ParseInt(string(count), 10, 64)
		if countInt == 0 {
			goto END
		}
		if countInt > 0 {
			var uid []interface{}
			uid = append(uid, userId)
			res, _ := redisCli.Sadd(userIdSetName, uid)
			if res == 1 {
				var order []interface{}
				order = append(order, userId)
				order = append(order, productTotal + 1 - countInt)
				redisCli.Hmset(productName, order)
				redisCli.Decr(countType)
				logger.Info("Order: userId=[%s],goodId=[%d]", order[0],order[1])
			}
		}
	}

	END:
	{
		logger.Info("Seckilling Done")
		fmt.Printf("Product ID: %d, done!\n", productId)
		PidFlag[productId] = false
	}
}
