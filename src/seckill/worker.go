package seckill

import (
	"fmt"
	"strconv"
	"helpers/iowrapper"
	"bytes"
	"encoding/binary"
)

func DealRequestQueue(productId int64, redisCli *iowrapper.RedisClient)  {
	productQueueName := "product_queue_" + strconv.FormatInt(productId, 10)
	fmt.Println(productQueueName)

	userId, _ := redisCli.Lpop(productQueueName)
	for userId != nil {
		count, _ := redisCli.Get("count");
		countInt := ConvertBytesToInt(count)
		if countInt > 0 {
			var uid []interface{}
			uid = append(uid, userId)
			status, _ := redisCli.Hget(productQueueName, uid)
			if status == 0 {
				var order []interface{}
				order = append(order, userId)
				order = append(order, 101 - count)
				redisCli.Hmset(productQueueName, order)
				redisCli.Incr("count")
			}

		}
		userId, _ = redisCli.Lpop(productQueueName)
	}
}

func ConvertBytesToInt(value []byte) int64 {
	buf :=  bytes .NewBuffer(value)
	var x int64
	binary.Read(buf, binary.BigEndian, &x)
	return x
}
