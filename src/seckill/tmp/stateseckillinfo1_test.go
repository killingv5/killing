package seckill

import (
	"fmt"
	"testing"
	"time"
)

func Test_redis(t *testing.T) {
	fmt.Println("start!")
	go ControlState1()
	time.Sleep(time.Second * 30)
	fmt.Println("finished!")
}
