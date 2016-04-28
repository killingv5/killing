package seckill
const(
	ERRNO_NONE = 0

	ERRNO_CONF_ERR      = 100000
	ERRNO_CONF_READFAIL = 100001

	ERRNO_REDIS_CONN_FAIL = 110001
	ERRNO_REDIS_SET_FAIL  = 110002
	ERRNO_REDIS_GET_FAIL  = 110003
	ERRNO_REDIS_DEL_FAIL  = 110004
	ERRNO_REDIS_RPUSH_FAIL = 11005


	ERRNO_SECKILL_FAIL            = 10000
	ERRNO_QUE_UERSECKILL_FAIL     = 10001
	ERRNO_QUE_PRODUCTSECKILL_FAIL = 10002

	ERRNO_PRODUCT_NOT_EXIST = 50000
	ERRNO_SECKILLING 	= 50001
	ERRNO_SECKILING_FAILED  = 50002

	ERRNO_PARSE_FAILED = 60000

)