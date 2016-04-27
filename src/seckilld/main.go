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
)

var(
	pidCountMap       map[int]int
	redisCli  		  *iowrapper.RedisClient
	conf              *seckill.Config
)

func init() {
	pidCountMap = make(map[int]int)
}

func seckillingHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if len(req.Form["userid"]) <= 0 || len(req.Form["productid"]) <= 0 {
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

	fmt.Println("before")
	err = seckill.Pushtoredis(req.Form["productid"][0], req.Form["userid"][0], redisCli)
	fmt.Println("after")
	if err != nil {
    	w.Write([]byte("unknow error"))
    	fmt.Println(err)
	} else {
    	w.Write([]byte("Hello"))
	}
}

func queryUserSeckillingInfoHandle(w http.ResponseWriter, req *http.Request) {
    req.ParseForm()
	if len(req.Form["userid"]) <= 0 || len(req.Form["productid"]) <= 0 {
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

	var retSt proSeckRet
	err, rets := seckill.QueryProductSeckingInfo(req.Form["productid"][0], redisCli)
	if err != nil {
		retSt = proSeckRet{112, make([]seckill.ProductSeckingInfo, 0)}
	} else {
		retSt = proSeckRet{0, rets}
	}
	fmt.Println(rets)

	//xxx := woqu{Error:1, List:[]ChatDb{ChatDb{12,15}, ChatDb{23, 343}}}
	retJson, err := json.Marshal(retSt)
	if err != nil {
		return
	}
	fmt.Println(string(retJson))

	w.Write([]byte(retJson))
}

func initFromConf() error {
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

	logfile := conf.GetValue("log","logfile")
	fmt.Println(logfile)
	_ = logger.SetupLogWithConf(logfile)
	/*
	if err != nil{
		return err
	}
	*/
	return nil
}

func initRedisCli(serverInfo string) error {
	fmt.Println(serverInfo)
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
	port := conf.GetValue("http","port")
    http.ListenAndServe(port, nil)
}

func main() {

	argc := len(os.Args)
	if (argc != 2){
		fmt.Println("usage: bin/seckill configFile")
		return
	}

	conf = seckill.SetConfig(os.Args[1])

	err := initFromConf()
	if err!=nil{
		return 
	}
	err = initWorker()
	if err != nil {
		fmt.Println(err)
		return
	}

	startHttpServer()
}
