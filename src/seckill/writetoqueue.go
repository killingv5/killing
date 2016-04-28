//write http request to redis queue
package seckill

import (
	logger "github.com/xlog4go"
	"helpers/iowrapper"
)

//push pid and uid of http request to a redis message queue
func PushToQueue(pid string, uid string, redisCli *iowrapper.RedisClient) (error) {
	key := PRODUCT_QUEUE + pid
	value := uid
	if redisCli != nil {
		_, err := redisCli.Rpush(key, value)
		if err != nil {
			logger.Error("error=[%s] key=[%s] value=[%s] err=[%s]", ERRNO_REDIS_RPUSH_FAIL,key, value, err.Error())
		}
		return err
	}

	return nil
}

//the main function for use
func PushToRedis(pid string, uid string, redisCli *iowrapper.RedisClient) ( error) {
	err := PushToQueue(pid, uid, redisCli)
	if err != nil {
		return err
	}
	
	return err
}
