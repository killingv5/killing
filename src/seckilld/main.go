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
	//"errors"
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



func paramCheck(req *http.Request, needUid bool, needSign bool) int {
	if len(req.Form["productid"]) <= 0 {
		return seckill.ERRNO_LACK_PROID
		//return errors.New("缺少商品信息!")
	}

	if needUid && len(req.Form["userid"]) <= 0 {
		return seckill.ERRNO_LACK_USRID
		//return errors.New("缺失用户信息!")
	}

	if !needSign {
		return seckill.ERRNO_NONE
	}

	if len(req.Form["sign"]) <= 0 {
		return seckill.ERRNO_LACK_SIGN
		//return errors.New("服务器未收到信息!")
	}

	var uidPid string
	if needUid {
		uidPid = req.Form["userid"][0] + req.Form["productid"][0]
	} else {
		uidPid = req.Form["productid"][0]
	}

	h := md5.New()
	h.Write([]byte(uidPid))

	if !strings.EqualFold(req.Form["sign"][0], hex.EncodeToString(h.Sum(nil))) {
		return seckill.ERRNO_SIGN_ERR
		//return errors.New("sign error")
	}

	return seckill.ERRNO_NONE
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
}

/**
* 查询商品列表
**/
func getProductListHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	str, err := seckill.GetProductList(redisCli)
	if err != nil {
    	w.Write([]byte("查询商品列表失败！"))
	} else {
    	w.Write([]byte(str))
	}
}

func seckillingHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	errno := seckill.ERRNO_SECKILLING
	defer func(){
		w.Write([]byte(seckill.MakeErrRet(errno)))
	}()

	//返回状态码,并复制给status
	errNo := paramCheck(req, true, needCheckSign)
	if errNo != seckill.ERRNO_NONE {
		errno = errNo
		//w.Write([]byte(err.Error()))
		return
	}

	pid, err := strconv.Atoi(req.Form["productid"][0])
	if err != nil {
		errno = seckill.ERRNO_PARA_NUM
		//w.Write([]byte("参数输入错误!"))
		return
	}

/*	value,okxx := seckill.PidFlag[int64(pid)]
	if okxx && !value {
		w.Write([]byte("活动结束,稍后请到查询页面查询结果"))
		return
	}*/

	_, ok := pidCountMap[pid]
	if !ok {
		errno = seckill.ERRNO_PRODUCT_NOT_EXIST
		//w.Write([]byte("商品信息不存在!"))
		return
	}

	err = seckill.PushToRedis(req.Form["productid"][0], req.Form["userid"][0], redisCli)
	errno = seckill.ERRNO_SECKILLING

}

func queryUserSeckillingInfoHandle(w http.ResponseWriter, req *http.Request) {

	retMap := make(map[string]int64)
	var status int64
	goodsId := int64(-1)
	errno := 0

	defer func(){
		if errno != 0 {
			w.Write([]byte(seckill.MakeErrRet(errno)))
		} else {
			retMap["errno"] = seckill.ERRNO_NONE
			retMap["status"] = status
			retMap["goodsid"] = goodsId
			retJson, _ := json.Marshal(retMap)
			w.Write([]byte(retJson))
			logger.Info("g_query_user_secking_info||timestamp=%s||ret=%s", time.Now().Format("2006-01-02 15:04:05"), retJson)
		}
	}()

	req.ParseForm()
	errNo := paramCheck(req, true, needCheckSign)
	if errNo != seckill.ERRNO_NONE {
		errno = errNo
		//w.Write([]byte(err.Error()))
		return
	}

	pid, err := strconv.Atoi(req.Form["productid"][0])
	if err != nil {
		errno = seckill.ERRNO_PARA_NUM
		//w.Write([]byte("参数输入错误!"))
		return
	}

	_, ok := pidCountMap[pid]
	if !ok {
		errno = seckill.ERRNO_PRODUCT_NOT_EXIST
		//w.Write([]byte("商品信息错误!"))
		return
	}

	info, err := seckill.QueryUserSeckillingInfo(req.Form["userid"][0], req.Form["productid"][0], redisCli)
	if err != nil {
		errno = seckill.ERRNO_QUE_UERSECKILL_FAIL
	} else {
		status = info.Status
		goodsId = info.Goodsid
	}

}

type proSeckRet struct {
	Error int                                `json:"error"`
	List  []seckill.ProductSeckingInfo       `json:"list"`
}

func queryProductSeckillingInfoHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	var retJson []byte
	errno := 0

	defer func(){
		if errno != 0 {
			w.Write([]byte(seckill.MakeErrRet(errno)))
		} else {
			w.Write([]byte(retJson))
		}
	}()

	//var status int64
	errNo := paramCheck(req, false, needCheckSign)
	if errNo != seckill.ERRNO_NONE {
		errno = errNo
		//w.Write([]byte(err.Error()))
		return
	}

	pid, err := strconv.Atoi(req.Form["productid"][0])
	if err != nil {
		errno = seckill.ERRNO_PARA_NUM
		//w.Write([]byte("参数输入错误!"))
		return
	}

	_, ok := pidCountMap[pid]
	if !ok {
		errno = seckill.ERRNO_PRODUCT_NOT_EXIST
		//w.Write([]byte("商品信息不存在!"))
		return
	}

	var retSt proSeckRet
	err, rets := seckill.QueryProductSeckingInfo(req.Form["productid"][0], redisCli)
	if err != nil {
		retSt = proSeckRet{seckill.ERRNO_QUE_PRODUCTSECKILL_FAIL, make([]seckill.ProductSeckingInfo, 0)}
	} else {
		retSt = proSeckRet{seckill.ERRNO_NONE, rets}
	}
	fmt.Println(rets)

	retJson, err = json.Marshal(retSt)
	if err != nil {
		errno = seckill.ERRNO_UNKNOW
		//w.Write([]byte("商品信息不存在!"))
		return
	}
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

	go startMisServer()

	startHttpServer()

}
