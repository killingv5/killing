package main

import (
    "net/http"
    "fmt"
    "strconv"
    "helpers/iowrapper"
)

var(
	pidCountMap       map[int]int
)

func init() {
	pidCountMap = make(map[int]int)	
}

func seckilling(w http.ResponseWriter, req *http.Request) {
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

func queryUserSeckillingInfo(w http.ResponseWriter, req *http.Request) {
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

func queryProductSeckillingInfo(w http.ResponseWriter, req *http.Request) {
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
	
	fmt.Println(pid)

    w.Write([]byte("Hello"))
}

func initWorker() {
	driverRedisClient := &iowrapper.RedisClient{
			Servers:        []string{"127.0.0.1:6379"},
		}
		err := driverRedisClient.Init()
		if err != nil {
			//logger.Error("init redis failed, err:%s", err.Error())
			fmt.Println("error")
			return
		}
}

func main() {

	initWorker()

    http.HandleFunc("/killing/seckilling", seckilling)
    http.HandleFunc("/killing/queryUserSeckillingInfo", queryUserSeckillingInfo)
    http.HandleFunc("/killing/queryProductSeckillingInfo", queryProductSeckillingInfo)
    http.ListenAndServe(":8001", nil)
}