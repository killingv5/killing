package seckill

import (
	"strconv"
	"encoding/json"
)

func Errno2Msg(errno int) string {
	msg, ok := Errmsg[errno]
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