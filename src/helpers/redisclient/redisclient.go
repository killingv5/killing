package redisclient

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

type RedisClient struct {

	// for pool
	pool         *redis.Pool
	redissvr     []string
	redissvrcnt  int
	conntimeout  int
	readtimeout  int
	writetimeout int
	maxidle      int
	maxactive    int

	// for rc
	expiresecond int
}

func NewRedisClient(redissvr map[string]string, conntimeout, readtimeout, writetimeout, maxidle, maxactive, expiresecond int) *RedisClient {

	rc := new(RedisClient)
	if rc == nil {
		return nil
	}

	redissvrarr := make([]string, 0)
	for _, v := range redissvr {
		redissvrarr = append(redissvrarr, v)
	}
	redissvrcnt := len(redissvrarr)

	rc.pool = GetRedisPool(redissvrarr, redissvrcnt, conntimeout, readtimeout, writetimeout, maxidle, maxactive)
	if rc.pool == nil {
		return nil
	}

	rc.redissvr = redissvrarr
	rc.redissvrcnt = redissvrcnt
	rc.conntimeout = conntimeout
	rc.readtimeout = readtimeout
	rc.writetimeout = writetimeout
	rc.maxidle = maxidle
	rc.maxactive = maxactive
	rc.expiresecond = expiresecond

	return rc
}

func (rc *RedisClient) Set(key, value string) error {
	c := rc.pool.Get()
	defer c.Close()

	reply, err := redis.String((c.Do("SET", key, value)))
	if err != nil {
		return err
	}

	// add redis key expire time.
	// ignore if error of expire command.
	rc.Expire(key, rc.expiresecond)

	if reply == "OK" {
		return nil
	} else {
		return errors.New("redisclient: unexpected reply of set")
	}
}

func (rc *RedisClient) Get(key string) (string, error) {
	c := rc.pool.Get()
	defer c.Close()

	reply, err := redis.String(c.Do("GET", key))
	if err != nil {
		return "", err
	}
	return reply, nil
}

func (rc *RedisClient) Setnx(key, value string) error {
	c := rc.pool.Get()
	defer c.Close()

	reply, err := redis.Int((c.Do("SETNX", key, value)))
	if err != nil {
		return err
	}

	// add redis key expire time.
	// ignore if error of expire command.
	rc.Expire(key, rc.expiresecond)

	if reply == 1 {
		return nil
	} else {
		return errors.New("redisclient: setnx fail of key exist")
	}
}

func (rc *RedisClient) Expire(key string, expiresecond int) error {
	c := rc.pool.Get()
	defer c.Close()

	reply, err := redis.Int(c.Do("EXPIRE", key, expiresecond))
	if err != nil {
		return err
	}

	if reply == 1 {
		return nil
	} else {
		return errors.New("redisclient: unexpected reply of expire")
	}
}

func (rc *RedisClient) Del(key string) error {
	c := rc.pool.Get()
	defer c.Close()

	_, err := redis.Int(c.Do("DEL", key))
	if err != nil {
		return err
	}
	//	if reply == 1 {
	//		return nil
	//	} else {
	// reply为0时说明key不存在
	//		return errors.New("redisclient: unexpected reply of del")
	//	}
	return nil
}

func (rc *RedisClient) ZRange(key string, start, stop int) ([]string, error) {
	c := rc.pool.Get()
	defer c.Close()

	reply, err := redis.Strings(c.Do("ZRANGE", key, start, stop))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (rc *RedisClient) ZRangeWithScores(key string, start, stop int) (map[string]string, error) {
	c := rc.pool.Get()
	defer c.Close()

	reply, err := redis.StringMap(c.Do("ZRANGE", key, start, stop, "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (rc *RedisClient) ZRangeByScore(key string, min, max int, minopen, maxopen bool) ([]string, error) {
	c := rc.pool.Get()
	defer c.Close()

	minstr := strconv.FormatInt(int64(min), 10)
	maxstr := strconv.FormatInt(int64(max), 10)
	if minopen {
		minstr = "(" + strconv.FormatInt(int64(min), 10)
	}

	if maxopen {
		maxstr = "(" + strconv.FormatInt(int64(max), 10)
	}

	reply, err := redis.Strings(c.Do("ZRANGEBYSCORE", key, minstr, maxstr))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (rc *RedisClient) ZRangeByScoreWithScores(key string, min, max int, minopen, maxopen bool) (map[string]string, error) {
	c := rc.pool.Get()
	defer c.Close()

	minstr := strconv.FormatInt(int64(min), 10)
	maxstr := strconv.FormatInt(int64(max), 10)
	if minopen {
		minstr = "(" + strconv.FormatInt(int64(min), 10)
	}

	if maxopen {
		maxstr = "(" + strconv.FormatInt(int64(max), 10)
	}

	reply, err := redis.StringMap(c.Do("ZRANGEBYSCORE", key, minstr, maxstr, "WITHSCORES"))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (rc *RedisClient) HGetall(key string) (map[string]string, error) {
	c := rc.pool.Get()
	defer c.Close()

	reply, err := redis.StringMap(c.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (rc *RedisClient) HGet(key, subkey string) (string, error) {
	c := rc.pool.Get()
	defer c.Close()

	reply, err := redis.String(c.Do("HGET", key, subkey))
	if err != nil {
		return "", err
	}

	return reply, nil
}

func (rc *RedisClient) HSet(key, subkey, value string) error {
	c := rc.pool.Get()
	defer c.Close()

	_, err := redis.Int(c.Do("HSET", key, subkey, value))
	if err != nil {
		return err
	}

	// add redis key expire time.
	// ignore if error of expire command.
	rc.Expire(key, rc.expiresecond)

	// no need to check reply of HSET
	// reply == 1 means HSET key subkey value, subkey not exist
	// reply == 0 means HSET key subkey value, subkey exists, but the value is already modified.
	/*
		if reply == 1 {
			return nil
		} else {
			return errors.New("redisclient: unexpected reply of hset")
		}
	*/

	return nil
}
