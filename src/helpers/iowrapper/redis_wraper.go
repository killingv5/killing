/*
   redis访问接口的包装，内部采取连接池实现
*/
package iowrapper

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	logger "github.com/xlog4go"
	"helpers/common"
	"time"
	"errors"
)

type RedisClient struct {
	Servers        []string
	ConnTimeoutMs  int
	WriteTimeoutMs int
	ReadTimeoutMs  int

	MaxIdle        int
	MaxActive      int
	IdleTimeoutS   int
	Password       string

	current_index  int
	pool           *redis.Pool
}

func (client *RedisClient) Close() {
	client.pool.Close()
}

func (client *RedisClient) Init() error {
	if len(client.Servers) == 0 {
		return fmt.Errorf("invalid Redis config servers:%s", client.Servers)
	}

	client.pool = &redis.Pool{
		MaxIdle:     client.MaxIdle,
		IdleTimeout: time.Duration(client.IdleTimeoutS) * time.Second,
		MaxActive:   client.MaxActive,
		Dial: func() (redis.Conn, error) {
			var c redis.Conn
			var err error
			for i := 0; i < len(client.Servers); i++ {
				//随机挑选一个IP
				index := common.RandIntn(len(client.Servers))
				client.current_index = index
				c, err = redis.DialTimeout("tcp", client.Servers[index],
					time.Duration(client.ConnTimeoutMs) * time.Millisecond,
					time.Duration(client.ReadTimeoutMs) * time.Millisecond,
					time.Duration(client.WriteTimeoutMs) * time.Millisecond)
				if err != nil {
					logger.Warn("warning=[redis_connect_failed] num=[%d] server=[%s] err=[%s]",
						i, client.Servers[index], err.Error())
				}
				//支持密码认证
				if len(client.Password) > 0 {
					if _, err_pass := c.Do("AUTH", client.Password); err_pass != nil {
						c.Close()
					}
				}
				if err == nil {
					logger.Info("info=[redis_connect_ok] num=[%d] server=[%s] err=[%s]",
						i, client.Servers[index])
					break
				}
			}
			return c, err
		},
	}

	return nil
}

func (client *RedisClient) Set(key string, value []byte) error {
	conn := client.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		logger.Error("error=[redis_set_failed] server=[%s] key=[%s] err=[%s]",
			client.Servers[client.current_index], key, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()
		_, err = conn_second.Do("SET", key, value)
		if err != nil {
			logger.Error("second error=[redis_set_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return err
		}
	}

	return nil
}

func (client *RedisClient) Get(key string) ([]byte, error) {
	conn := client.pool.Get()
	defer conn.Close()

	value, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			logger.Info("error=[redis_get_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return nil, err
		} else {
			logger.Error("error=[redis_get_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
		}

		conn_second := client.pool.Get()
		defer conn_second.Close()

		value, err = redis.Bytes(conn_second.Do("GET", key))
		if err != nil {
			if err.Error() == "redigo: nil returned" {
				logger.Info("second error=[redis_get_failed] server=[%s] key=[%s] err=[%s]",
					client.Servers[client.current_index], key, err.Error())
			} else {
				logger.Error("second error=[redis_get_failed] server=[%s] key=[%s] err=[%s]",
					client.Servers[client.current_index], key, err.Error())
			}
			return nil, err
		}
	}

	return value, nil
}

func (client *RedisClient) Rpush(key string, value string) (int64, error) {
	conn := client.pool.Get()
	defer conn.Close()

	list_len, err := redis.Int64(conn.Do("RPUSH", key, value))
	if err != nil {
		logger.Error("error=[redis_rpush_failed] server=[%s] key=[%s] err=[%s]",
			client.Servers[client.current_index], key, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()

		list_len, err = redis.Int64(conn_second.Do("RPUSH", key, value))
		if err != nil {
			logger.Error("second error=[redis_rpush_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return -1, err
		}
	}

	return list_len, nil
}

func (client *RedisClient) Lpop(key string) (string, error) {
	conn := client.pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("LPOP", key))
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			logger.Info("error=[redis_lpop_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return "", err
		} else {
			logger.Error("error=[redis_lpop_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
		}

		conn_second := client.pool.Get()
		defer conn_second.Close()

		value, err = redis.String(conn_second.Do("LPOP", key))
		if err != nil {
			if err.Error() == "redigo: nil returned" {
				logger.Info("second error=[redis_lpop_failed] server=[%s] key=[%s] err=[%s]",
					client.Servers[client.current_index], key, err.Error())
			} else {
				logger.Error("second error=[redis_lpop_failed] server=[%s] key=[%s] err=[%s]",
					client.Servers[client.current_index], key, err.Error())
			}
			return "", err
		}
	}

	return value, nil
}

func (client *RedisClient) Llen(key string) (int64, error) {
	conn := client.pool.Get()
	defer conn.Close()

	value, err := redis.Int64(conn.Do("LLEN", key))
	if err != nil {
		logger.Error("error=[redis_llen_failed] server=[%s] key=[%s] err=[%s]",
			client.Servers[client.current_index], key, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()

		value, err = redis.Int64(conn_second.Do("LLEN", key))
		if err != nil {
			logger.Error("second error=[redis_llen_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return -1, err
		}
	}
	return value, nil
}

func (client *RedisClient) Del(keys []interface{}) (int64, error) {
	conn := client.pool.Get()
	defer conn.Close()

	value, err := redis.Int64(conn.Do("DEL", keys...))
	if err != nil {
		logger.Error("error=[redis_del_failed] server=[%s] keys=[%v] err=[%s]",
			client.Servers[client.current_index], keys, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()

		value, err = redis.Int64(conn_second.Do("DEL", keys...))
		if err != nil {
			logger.Error("second error=[redis_del_failed] server=[%s] keys=[%v] err=[%s]",
				client.Servers[client.current_index], keys, err.Error())
			return -1, err
		}
	}

	return value, nil
}

func (client *RedisClient) Hmset(key string, value []interface{}) (string, error) {
	conn := client.pool.Get()
	defer conn.Close()

	//input_params := make([]interface{}, len(value)+1)
	var input_params []interface{}
	input_params = append(input_params, key)
	input_params = append(input_params, value...)
	res, err := redis.String(conn.Do("HMSET", input_params...))
	if err != nil {
		logger.Error("error=[redis_hmset_failed] server=[%s] key=[%s] err=[%s]",
			client.Servers[client.current_index], key, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()

		res, err = redis.String(conn_second.Do("HMSET", input_params...))
		if err != nil {
			logger.Error("second error=[redis_hmset_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return "", err
		}
	}

	return res, nil
}

func (client *RedisClient) Hdel(key string, value []interface{}) (int64, error) {
	conn := client.pool.Get()
	defer conn.Close()

	//input_params := make([]interface{}, len(value)+1)
	var input_params []interface{}
	input_params = append(input_params, key)
	input_params = append(input_params, value...)
	res, err := redis.Int64(conn.Do("HDEL", input_params...))
	if err != nil {
		logger.Error("error=[redis_hdel_failed] server=[%s] key=[%s] err=[%s]",
			client.Servers[client.current_index], key, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()

		res, err = redis.Int64(conn_second.Do("HDEL", input_params...))
		if err != nil {
			logger.Error("second error=[redis_hdel_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return -1, err
		}
	}

	return res, nil
}

func (client *RedisClient) Hkeys(key string) ([]string, error) {
	conn := client.pool.Get()
	defer conn.Close()

	res, err := redis.Strings(conn.Do("HKEYS", key))
	if err != nil {
		logger.Error("error=[redis_hkeys_failed] server=[%s] key=[%s] err=[%s]",
			client.Servers[client.current_index], key, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()

		res, err = redis.Strings(conn_second.Do("HKEYS", key))
		if err != nil {
			logger.Error("second error=[redis_hkeys_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return nil, err
		}
	}

	return res, nil
}

func (client *RedisClient) Keys(key string) ([]string, error) {
	conn := client.pool.Get()
	defer conn.Close()

	res, err := redis.Strings(conn.Do("KEYS", key))
	if err != nil {
		logger.Error("error=[redis_keys_failed] server=[%s] key=[%s] err=[%s]",
			client.Servers[client.current_index], key, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()

		res, err = redis.Strings(conn_second.Do("KEYS", key))
		if err != nil {
			logger.Error("second error=[redis_keys_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return nil, err
		}
	}

	return res, nil
}

func (client *RedisClient) Hgetall(key string) ([]string, error) {
	conn := client.pool.Get()
	defer conn.Close()

	res, err := redis.Strings(conn.Do("HGETALL", key))
	if err != nil {
		logger.Error("error=[redis_hgetall_failed] server=[%s] key=[%s] err=[%s]",
			client.Servers[client.current_index], key, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()

		res, err = redis.Strings(conn_second.Do("HGETALL", key))
		if err != nil {
			logger.Error("second error=[redis_hgetall_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return nil, err
		}
	}

	return res, nil
}

func (client *RedisClient) Hget(key string, field string) (string, error) {
	conn := client.pool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("HGET", key, field))

	if err != nil {
		logger.Error("error=[redis_hget_failed] server=[%s] key=[%s] err=[%s]",
			client.Servers[client.current_index], key, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()

		res, err = redis.String(conn_second.Do("HGET", key, field))
		if err != nil {
			logger.Error("second error=[redis_hget_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return "", err
		}
	}

	return res, nil
}

func (client *RedisClient) Sadd(key string, value []interface{}) (int64, error) {
	conn := client.pool.Get()
	defer conn.Close()

	var input_params []interface{}
	input_params = append(input_params, key)
	input_params = append(input_params, value...)
	res, err := redis.Int64(conn.Do("SADD", input_params...))
	if err != nil {
		logger.Error("error=[redis_sadd_failed] server=[%s] key=[%s] err=[%s]",
			client.Servers[client.current_index], key, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()

		res, err = redis.Int64(conn_second.Do("SADD", input_params...))
		if err != nil {
			logger.Error("second error=[redis_sadd_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return -1, err
		}
	}

	return res, nil
}

func (client *RedisClient) Smembers(key string) ([]string, error) {
	conn := client.pool.Get()
	defer conn.Close()

	res, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		logger.Error("error=[redis_smembers_failed] server=[%s] key=[%s] err=[%s]",
			client.Servers[client.current_index], key, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()

		res, err = redis.Strings(conn_second.Do("SMEMBERS", key))
		if err != nil {
			logger.Error("second error=[redis_smembers_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return nil, err
		}
	}

	return res, nil
}

func (client *RedisClient) Decr(key string) (int64, error) {
	conn := client.pool.Get()
	defer conn.Close()

	res, err := redis.Int64(conn.Do("DECR", key))
	if err != nil {
		logger.Error("error=[redis_hget_failed] server=[%s] key=[%s] err=[%s]",
			client.Servers[client.current_index], key, err.Error())

		conn_second := client.pool.Get()
		defer conn_second.Close()

		res, err = redis.Int64(conn.Do("DECR", key))
		if err != nil {
			logger.Error("second error=[redis_hget_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return 0, err
		}
	}
	return res, err
}

func (client *RedisClient) BLpop(key string, timeout int64) (string, error) {
	conn := client.pool.Get()
	value := ""
	defer conn.Close()

	rslt,err := conn.Do("BLPOP", key, timeout)
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			logger.Info("error=[redis_lpop_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
			return "", err
		} else {
			logger.Error("error=[redis_lpop_failed] server=[%s] key=[%s] err=[%s]",
				client.Servers[client.current_index], key, err.Error())
		}

		conn_second := client.pool.Get()
		defer conn_second.Close()

		rslt,err = conn.Do("BLPOP", key, timeout)
		if err != nil {
			if err.Error() == "redigo: nil returned" {
				logger.Info("second error=[redis_lpop_failed] server=[%s] key=[%s] err=[%s]",
					client.Servers[client.current_index], key, err.Error())
			} else {
				logger.Error("second error=[redis_lpop_failed] server=[%s] key=[%s] err=[%s]",
					client.Servers[client.current_index], key, err.Error())
			}
			return "", err
		}
	}

	if val, ok := rslt.([]interface{});ok {
		if len(val) < 2 {
			return "", errors.New("Redis return err")
		}
		if valbit ,ok:= val[1].([]byte);ok {
			value =  string(valbit)
		}
	}

	return value, nil
}
