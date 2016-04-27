package seckill

import (
	"fmt"
	"strconv"
	"helpers/iowrapper"
)

func DealRequestQueue(productId int64, redisCli *iowrapper.RedisClient)  {
	productQueueName := "product_queue_" + strconv.FormatInt(productId, 10)
	productName := "product_" + strconv.FormatInt(productId, 10)

	userId, _ := redisCli.Lpop(productQueueName)
	for userId != "" {
		count, _ := redisCli.Get("count");
		countInt, _ := strconv.ParseInt(string(count), 10, 64)
		if countInt > 0 {
			status, _ := redisCli.Hget(productName, userId)
			statusInt, _ := strconv.ParseInt(status, 10, 64)
			if statusInt == 0 {
				var order []interface{}
				order = append(order, userId)
				order = append(order, 101 - countInt)
				redisCli.Hmset(productName, order)
				fmt.Println(countInt)
				redisCli.Decr("count")
			}

		}
		userId, _ = redisCli.Lpop(productQueueName)
	}
}
