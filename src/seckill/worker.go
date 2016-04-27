package seckill

import (
	"strconv"
	"helpers/iowrapper"
	logger "github.com/xlog4go"
)

func DealRequestQueue(productId int64, redisCli *iowrapper.RedisClient)  {
	productType := strconv.FormatInt(productId, 10)
	countType := "count_" + productType
	productQueueName := "product_queue_" + productType
	userIdSetName := "userid_set_" + productType
	productName := "product_" + productType

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
				order = append(order, 101 - countInt)
				redisCli.Hmset(productName, order)
				redisCli.Decr(countType)
				logger.Info("Order: %s", order)
			}
		}
	}

	END:
	{
		logger.Info("Seckilling Done")
	}
	
}
