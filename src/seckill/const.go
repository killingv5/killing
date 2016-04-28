package seckill

import (
	"strconv"
	"encoding/json"
)
const (
	//stuct in redis
	PRODUCT_QUEUE = "product_queue_"
	COUNT_TYPE = "count_"
	USERID_SET = "userid_set_"
	PRODUCT_HASH = "product_hash_"

	//state code
	STATE_NOT_STARTED = 10 //抢单未开始
	STATE_ING = 11 //抢单进行中
	STATE_ENDED = 12 //抢单结束
	STATE_NOT_EXIST = 13 //商品不存在

	//status code
	SECKILLING_NOT_START = 0
	SECKILLING_SUCCESS = 1
	SECKILLING_FAIL = 2
	PRODUCT_NOT_EXIST = 3

	PRODUCTINFO = "product_info"
	TIMEFORMAT  = "2006-01-02 15:04:05"

	//error code
	ERRNO_NONE = 0

	ERRNO_CONF_ERR      = 100000
	ERRNO_CONF_READFAIL = 100001

	ERRNO_REDIS_CONN_FAIL = 110001
	ERRNO_REDIS_SET_FAIL  = 110002
	ERRNO_REDIS_GET_FAIL  = 110003
	ERRNO_REDIS_DEL_FAIL  = 110004
	ERRNO_REDIS_RPUSH_FAIL = 110005

	ERRNO_SECKILL_FAIL            = 10000
	ERRNO_QUE_UERSECKILL_FAIL     = 10001
	ERRNO_QUE_PRODUCTSECKILL_FAIL = 10002

	ERRNO_SECKILLING 	    = 50000
	ERRNO_SECKILING_FAILED  = 50001
	ERRNO_PRODUCT_NOT_EXIST = 50002
	ERRNO_LACK_PROID        = 50003
	ERRNO_LACK_USRID        = 50004
	ERRNO_LACK_SIGN         = 50005
	ERRNO_SIGN_ERR          = 50006

	ERRNO_PARA_NUM          = 60000
	ERRNO_PARSE_FAILED      = 60001

	ERRNO_UNKNOW            = 70001
)

var errmsg map[int]string

func init() {
	errmsg = make(map[int]string)
	errmsg[ERRNO_NONE] = "成功"
	errmsg[ERRNO_SECKILL_FAIL] = "查询失败"
	errmsg[ERRNO_QUE_UERSECKILL_FAIL] = "查询失败"
	errmsg[ERRNO_QUE_PRODUCTSECKILL_FAIL] = "查询失败"
	errmsg[ERRNO_SECKILLING] = "正在秒杀"
	errmsg[ERRNO_SECKILING_FAILED] = "秒杀失败"
	errmsg[ERRNO_PRODUCT_NOT_EXIST] = "商品不存在"
	errmsg[ERRNO_LACK_PROID] = "缺少商品ID"
	errmsg[ERRNO_LACK_USRID] = "缺少用户ID"
	errmsg[ERRNO_LACK_SIGN] = "缺少校验值"
	errmsg[ERRNO_SIGN_ERR] = "信号校验错误"
	errmsg[ERRNO_PARA_NUM] = "参数不合法"
	errmsg[ERRNO_UNKNOW] = "unknow error"
}

func Errno2Msg(errno int) string {
	msg, ok := errmsg[errno]
	if !ok {
		return "unknow error"
	} else {
		return msg
	}
}

func MakeErrRet(err int) string {
	retMap := make(map[string]string)
	retMap["errno"] = strconv.Itoa(err)
	retMap["msg"] = Errno2Msg(err)
	retJson, _ := json.Marshal(retMap)
	return string(retJson)
}
