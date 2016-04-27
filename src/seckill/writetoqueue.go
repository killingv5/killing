//write http request to redis queue
package seckill

import (
	//"fmt"
	logger "github.com/xlog4go"
	"helpers/iowrapper"
)

const (
	PRODUCTQUEUE = "product_queue_"
	PRODUCT      = "product_"
)

//push pid, uid and initial status to a redis hash map
func pushtohash(pid string, uid string, redisCli *iowrapper.RedisClient) (error) {
	var values []interface{}
	key := PRODUCT + pid

	status := "0"
	values = append(values, uid)
	values = append(values, status)
	if redisCli != nil {
	    _, err := redisCli.Hmset(key, values)
	    if err != nil {
            logger.Error("error=[redis_hmset_hash_failed] key=[%s] uid=[%s] status=[%s] err=[%s]", key, uid, status, err.Error())
	}
	return err
	}


	return nil
}

//push pid and uid of http request to a redis message queue
func pushtoqueue(pid string, uid string, redisCli *iowrapper.RedisClient) (error) {
	key := PRODUCTQUEUE + pid
	value := uid
	if redisCli != nil {
		_, err := redisCli.Rpush(key, value)
	if err != nil {
		logger.Error("error=[redis_rpush_failed] key=[%s] value=[%s] err=[%s]", key, value, err.Error())
	}
	return err
	}


	return nil
}

//the main function for use
func Pushtoredis(pid string, uid string, redisCli *iowrapper.RedisClient) ( error) {
	err := pushtohash(pid, uid, redisCli)
	if err != nil {
		return err
	}

	err = pushtoqueue(pid, uid, redisCli)
	if err != nil {
		return err
	}
	return err
}
