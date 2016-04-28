package main

import (
	"net/http"
	"fmt"
	"strconv"
	"helpers/iowrapper"
	"encoding/json"
	"seckill"
	"os"
	logger "github.com/xlog4go"
	"time"
	"errors"
	"crypto/md5"
	"encoding/hex"
	"strings"
)

var (
	pidCountMap       map[int]int
	redisCli                  *iowrapper.RedisClient
	serverInfo string
	logFile string
	needCheckSign bool
)

func init() {
	pidCountMap = make(map[int]int)
	needCheckSign = false
}

func paramCheck(req *http.Request, needUid bool, needSign bool) error {
	if len(req.Form["productid"]) <= 0 {
		return errors.New("productid miss")
	}

	if needUid && len(req.Form["userid"]) <= 0 {
		return errors.New("userid miss")
	}

	if !needSign {
		return nil
	}

	if len(req.Form["sign"]) <= 0 {
		return errors.New("sign miss")
	}

	var uidpid string
	if needUid {
		uidpid = req.Form["userid"][0] + req.Form["productid"][0]
	} else {
		uidpid = req.Form["productid"][0]
	}

	h := md5.New()
	h.Write([]byte(uidpid))
	if !strings.EqualFold(req.Form["sign"][0], hex.EncodeToString(h.Sum(nil))) {
		return errors.New("sign error")
	}

	return nil
}

/**
* 清空数据库
**/
func flushHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if len(req.Form["productid"]) <= 0 {
		w.Write([]byte("参数输入错误!"))
		return
	}

	pid, err := strconv.Atoi(req.Form["productid"][0])
	if err != nil {
		w.Write([]byte("参数输入错误!"))
		return
	}

	_, ok := pidCountMap[pid]
	if !ok {
		w.Write([]byte("商品信息不存在!"))
		return
	}

	err = seckill.CleanProduct(req.Form["productid"][0], redisCli)
	if err != nil {
		w.Write([]byte("数据清空失败！"))
		//fmt.Println(err)
	} else {
		w.Write([]byte("数据清空成功！"))
	}
}

/**
* 添加商品
**/
func addProductHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if len(req.Form["productid"]) <= 0 || len(req.Form["productnum"]) <=0 || len(req.Form["starttime"]) <=0 {
		w.Write([]byte("param error !"))
		return
	}
	seckill.AddProduct(req.Form["productid"][0], req.Form["productnum"][0], req.Form["starttime"][0], redisCli)
	w.Write([]byte("商品添加成功！"))
	pid, err := strconv.Atoi(req.Form["productid"][0])
	if err != nil {
		w.Write([]byte("参数输入错误!"))
		return
	}

	count, err := strconv.Atoi(req.Form["productnum"][0])
	if err != nil {
		w.Write([]byte("参数输入错误!"))
		return
	}
	pidCountMap[pid] = count
	go seckill.DealRequestQueue(int64(pid), redisCli)
}

/**
* 查询商品列表
**/
func getProductListHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	str, err := seckill.GetProductList(redisCli)
	if err != nil {
    	w.Write([]byte("查询商品列表失败！"))
    	fmt.Println(err)
	} else {
    	w.Write([]byte(str))
	}
}

func seckillingHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	err := paramCheck(req, true, needCheckSign)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	_ , err = strconv.Atoi(req.Form["productid"][0])
	if err != nil {
		w.Write([]byte("参数输入错误!"))
		return
	}

	state := seckill.GetPidState(req.Form["productid"][0])
	switch state {
	case seckill.STATE_NOT_STARTED:
		w.Write([]byte("秒杀未开始!"))
		return
	case seckill.STATE_ENDED:
		w.Write([]byte("秒杀已结束!"))
		return
	case seckill.STATE_NOT_EXIST:
		w.Write([]byte("商品信息错误!"))	
		return
	}

	err = seckill.PushToRedis(req.Form["productid"][0], req.Form["userid"][0], redisCli)
/*	if err != nil {
		//seckill.PushToRedis方法中打印log
		w.Write([]byte("unknow error"))
		//logger.Error("errno=[%s],err=[%s]", seckill.ERRNO_SECKILL_FAIL, err.Error())
	} */
	w.Write([]byte("排队中，结果请稍后查询..."))


	/*retJson, err := json.Marshal(retMap)
	if err != nil {
		w.Write([]byte("Unknow error !"))
		logger.Error("retMap to retJson falied, err:%s",err.Error())
		return
	}
	w.Write([]byte(retJson))
	//fmt.Println(err)*/
}

func queryUserSeckillingInfoHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	err := paramCheck(req, true, needCheckSign)
	if err != nil {
		w.Write([]byte("参数输入错误!"))
		return
	}

	_ , err = strconv.Atoi(req.Form["productid"][0])
	if err != nil {
		w.Write([]byte("参数输入错误!"))
		return
	}

	state := seckill.GetPidState(req.Form["productid"][0])
	if state == seckill.STATE_NOT_EXIST {
		w.Write([]byte("商品信息错误!"))	
		return
	}

	/*

	_, ok := pidCountMap[pid]
	if !ok {
		w.Write([]byte("商品信息错误!"))
		return
	}*/

	retMap := make(map[string]int64)
	info, err := seckill.QueryUserSeckillingInfo(req.Form["userid"][0], req.Form["productid"][0], redisCli)
	if err != nil {
		//retMap["errno"] = seckill.ERRNO_QUE_UERSECKILL_FAIL
		w.Write([]byte("很遗憾,没有秒杀到 ~"))
		//logger.Error("errno=[%s], err=[%s]", seckill.ERRNO_QUE_UERSECKILL_FAIL, err.Error())
	} else {
		retMap["errno"] = seckill.ERRNO_NONE
		retMap["status"] = info.Status
		retMap["goodsid"] = info.Goodsid

	}

	retJson, err := json.Marshal(retMap)
	if err != nil {
		w.Write([]byte("很遗憾,没有秒杀到 ~"))
		return
	}
	w.Write([]byte(retJson))
	logger.Info("g_query_user_secking_info||timestamp=%s||ret=%s", time.Now().Format("2006-01-02 15:04:05"), retJson)
}

type proSeckRet struct {
	Error int                                `json:"error"`
	List  []seckill.ProductSeckingInfo       `json:"list"`
}

func queryProductSeckillingInfoHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	err := paramCheck(req, false, needCheckSign)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	pid, err := strconv.Atoi(req.Form["productid"][0])
	if err != nil {
		w.Write([]byte("参数输入错误!"))
		return
	}

	_, ok := pidCountMap[pid]
	if !ok {
		w.Write([]byte("商品信息不存在!"))
		return
	}

	var retSt proSeckRet
	err, rets := seckill.QueryProductSeckingInfo(req.Form["productid"][0], redisCli)
	if err != nil {
		retSt = proSeckRet{seckill.ERRNO_QUE_PRODUCTSECKILL_FAIL, make([]seckill.ProductSeckingInfo, 0)}
	} else {
		retSt = proSeckRet{seckill.ERRNO_NONE, rets}
	}
	// fmt.Println(rets)

	retJson, err := json.Marshal(retSt)
	if err != nil {
		w.Write([]byte("商品信息不存在!"))
		return
	}
	w.Write([]byte(retJson))
	logger.Info("g_query_product_secking_info||timestamp=%s||ret=%s", time.Now().Format("2006-01-02 15:04:05"), retJson)

}

func initFromConf(configFile string) error {
	conf := seckill.SetConfig(configFile)
	serverInfo = conf.GetValue("redis", "serverInfo")
	//fmt.Println(serverInfo)
	//init logger
	logFile = conf.GetValue("log", "logfile")

	productId := conf.GetValue("product", "productid")
	productNum := conf.GetValue("product", "productnum")
	productid, _ := strconv.Atoi(productId);
	productnum, _ := strconv.Atoi(productNum);
	//fmt.Println(productid)
	//fmt.Println(productnum)
	pidCountMap[productid] = productnum

	needcheck := conf.GetValue("sign", "needcheck")
	if strings.EqualFold("1", needcheck) {
		needCheckSign = true
	}

	return nil
}

func initRedisCli(serverInfo string) error {
	redisCli = &iowrapper.RedisClient{
		Servers:        []string{serverInfo},
	}

	err := redisCli.Init()
	return err
}

func initWorker() error {
	for k, _ := range pidCountMap {
		go seckill.DealRequestQueue(int64(k), redisCli)
		fmt.Println(k)
	}
	return nil
}

func initController() error {
	if redisCli == nil {
		return errors.New("controller init failed : redis init failed")
	}
	go seckill.ControlState(redisCli)
	return nil
}

func startHttpServer() {
	http.HandleFunc("/killing/seckilling", seckillingHandle)
	http.HandleFunc("/killing/queryUserSeckillingInfo", queryUserSeckillingInfoHandle)
	http.HandleFunc("/killing/queryProductSeckillingInfo", queryProductSeckillingInfoHandle)
	http.ListenAndServe(":8001", nil)
}

func startMisServer() {
	http.HandleFunc("/killing/cleandb", flushHandle)
	http.HandleFunc("/killing/addproduct", addProductHandle)
	http.HandleFunc("/killing/getproductlist", getProductListHandle)
	http.ListenAndServe(":9001", nil)
}

func main() {

	argc := len(os.Args)
	if (argc != 2) {
		fmt.Println("usage bin/seckill configFile")
		return
	}

	err := initFromConf(os.Args[1])
	if err != nil {
		fmt.Println("init config failed,err:%s", err.Error())
		return
	}

	err = logger.SetupLogWithConf(logFile)
	if err != nil {
		fmt.Println("init log fail: %s", err.Error())
		return
	}
	defer logger.Close()

	err = initRedisCli(serverInfo)
	if err != nil {
		logger.Error("init redis failed,err:%s", err.Error())
		return
	}

	err = initWorker()
	if err != nil {
		//fmt.Println(err)
		logger.Error("init worker failed,err:%s", err.Error())
		return
	}

	err = initController()
	if err != nil {
		logger.Error("init controller failed,err:%s", err.Error())
		return
	}

	go startMisServer()

	startHttpServer()

}
