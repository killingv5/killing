package main

import (
    "net/http"
    "fmt"
    "strconv"
    "helpers/iowrapper"
	"encoding/json"
)

var(
	pidCountMap       map[int]int
	redisCli  		  *iowrapper.RedisClient
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

	uid, err := strconv.Atoi(req.Form["userid"][0])
	if err != nil {
		w.Write([]byte("param error !"))
		return
	}

	_, ok := pidCountMap[pid]
	if !ok {
		w.Write([]byte("no productid !"))
		return
	}

	fmt.Println(pid)
	fmt.Println(uid)

    w.Write([]byte("Hello"))
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

	uid, err := strconv.Atoi(req.Form["userid"][0])
	if err != nil {
		w.Write([]byte("param error !"))
		return
	}

	_, ok := pidCountMap[pid]
	if !ok {
		w.Write([]byte("no productid !"))
		return
	}

	retMap := make(map[string]int)
	retMap["errno"] = 0
	retMap["status"] = 1
	retMap["goodsid"] = 12

	retJson, err := json.Marshal(retMap)
	if err != nil {
		w.Write([]byte("unknow error !"))
		return
	}
	w.Write([]byte(retJson))

	fmt.Println(pid)
	fmt.Println(uid)
}

type xxx struct {
	Userid          int `json:"userid"`
	Goodsid         int `json:"goodsid"`
}

type woqu struct {
	Error int  `json:"error"`
	List []ChatDb `json:"list"`
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

	xxx := woqu{Error:1, List:[]ChatDb{ChatDb{12,15}, ChatDb{23, 343}}}
	retJson, err := json.Marshal(xxx)
	if err != nil {
		return
	}
	fmt.Println(string(retJson))

	w.Write([]byte(retJson))
}

func initRedisCli() error {
	redisCli := &iowrapper.RedisClient{
			Servers:        []string{"127.0.0.1:6379"},
	}

	err := redisCli.Init()
	return err
}

func initWorker() error{
	// start worker
	for k, _ := range pidCountMap {
		// go xxxWorker_fun(k, redisCli)
	}
	return nil
}

func startHttpServer() {
	http.HandleFunc("/killing/seckilling", seckillingHandle)
    http.HandleFunc("/killing/queryUserSeckillingInfo", queryUserSeckillingInfoHandle)
    http.HandleFunc("/killing/queryProductSeckillingInfo", queryProductSeckillingInfoHandle)
    http.ListenAndServe(":8001", nil)
}

func main() {

	err := initRedisCli()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = initWorker()
	if err != nil {
		fmt.Println(err)
		return
	}

	startHttpServer()
}