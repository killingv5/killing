package seckill

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

)
