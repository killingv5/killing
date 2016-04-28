package main

import (
    "net/http"
    "fmt"
    "strconv"
    "helpers/iowrapper"
	"encoding/json"
	"seckill"
	"os"
	//logger "github.com/xlog4go"
	"errors"
	"crypto/md5"
    "encoding/hex"
    "strings"
)

var(
	pidCountMap       map[int]int
	redisCli  		  *iowrapper.RedisClient
	needCheckSign	  bool
)

func init() {
	pidCountMap = make(map[int]int)
	needCheckSign = false
}

func paramCheck(req *http.Request, needUid bool, needSign bool) error{ 
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

func flushHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if len(req.Form["productid"]) <= 0 {
		w.Write([]byte("param error !"))
		return
	}

	pid, err := strconv.Atoi(req.Form["productid"][0])
	if err != nil {
		w.Write([]byte("param error !"))
		return
	}

	_, ok := pidCountMap[pid]
	if !ok {
		w.Write([]byte("no productid !"))
		return
	}

	err = seckill.CleanProduct(req.Form["productid"][0], redisCli)
	if err != nil {
    	w.Write([]byte("数据清空失败！"))
    	fmt.Println(err)
	} else {
    	w.Write([]byte("数据清空成功！"))
	}
}

func seckillingHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()

    err := paramCheck(req, true, needCheckSign)
    if err != nil {
    	w.Write([]byte(err.Error()))
    	return
    }

	pid, err := strconv.Atoi(req.Form["productid"][0])
	if err != nil {
		w.Write([]byte("param error !"))
		return
	}

	_, ok := pidCountMap[pid]
	if !ok {
		w.Write([]byte("no productid !"))
		return
	}

	err = seckill.PushToRedis(req.Form["productid"][0], req.Form["userid"][0], redisCli)
	if err != nil {
    	w.Write([]byte("unknow error"))
    	fmt.Println(err)
	} else {
    	w.Write([]byte("排队中，结果请稍后查询"))
	}
}

func queryUserSeckillingInfoHandle(w http.ResponseWriter, req *http.Request) {
    req.ParseForm()
	err := paramCheck(req, true, needCheckSign)
    if err != nil {
    	w.Write([]byte(err.Error()))
    	return
    }

	pid, err := strconv.Atoi(req.Form["productid"][0])
	if err != nil {
		w.Write([]byte("param error !"))
		return
	}

	_, ok := pidCountMap[pid]
	if !ok {
		w.Write([]byte("no productid !"))
		return
	}

	retMap := make(map[string]int64)
	info, err := seckill.QueryUserSeckillingInfo(req.Form["userid"][0], req.Form["productid"][0], redisCli)
	if err != nil {
		retMap["errno"] = 1001
	} else {
		retMap["errno"] = 0
		retMap["status"] = info.Status
		retMap["goodsid"] = info.Goodsid
		
	}

	retJson, err := json.Marshal(retMap)
	if err != nil {
		w.Write([]byte("unknow error !"))
		return
	}
	w.Write([]byte(retJson))
}

type proSeckRet struct {
	Error int  							`json:"error"`
	List []seckill.ProductSeckingInfo 	`json:"list"`
}

func queryProductSeckillingInfoHandle(w http.ResponseWriter, req *http.Request) {
	err := paramCheck(req, false, needCheckSign)
    if err != nil {
    	w.Write([]byte(err.Error()))
    	return
    }

	pid, err := strconv.Atoi(req.Form["productid"][0])
	if err != nil {
		w.Write([]byte("param error !"))
		return
	}

	_, ok := pidCountMap[pid]
	if !ok {
		w.Write([]byte("no productid !"))
		return
	}

	var retSt proSeckRet
	err, rets := seckill.QueryProductSeckingInfo(req.Form["productid"][0], redisCli)
	if err != nil {
		retSt = proSeckRet{112, make([]seckill.ProductSeckingInfo, 0)}
	} else {
		retSt = proSeckRet{0, rets}
	}
	fmt.Println(rets)

	retJson, err := json.Marshal(retSt)
	if err != nil {
		return
	}
	fmt.Println(string(retJson))

	w.Write([]byte(retJson))
}

func initFromConf(configFile string) error {
	conf := seckill.SetConfig(configFile)
	serverInfo := conf.GetValue("redis","serverInfo")
	fmt.Println(serverInfo)
	if err := initRedisCli(serverInfo);err != nil{
		return err
	}

	productId   := conf.GetValue("product","productid")
	productNum  := conf.GetValue("product","productnum")
	productid, _ := strconv.Atoi(productId);
	productnum, _ := strconv.Atoi(productNum);
	fmt.Println(productid)
	fmt.Println(productnum)
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
	//	Servers:        []string{"127.0.0.1:6379"},
	}

	err := redisCli.Init()
	return err
}

func initWorker() error{
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
    http.ListenAndServe(":9001", nil)
}

func main() {

	argc := len(os.Args)
	if (argc != 2){
		fmt.Println("usage bin/seckill configFile")
		return
	}

	err := initFromConf(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	err = initWorker()
	if err != nil {
		fmt.Println(err)
		return
	}
	go startMisServer()
	
	startHttpServer()
	
}
