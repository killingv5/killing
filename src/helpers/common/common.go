package common

import (
	"math"
	"math/rand"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func CallerName() string {
	var pc uintptr
	var file string
	var line int
	var ok bool
	if pc, file, line, ok = runtime.Caller(1); !ok {
		return ""
	}
	name := runtime.FuncForPC(pc).Name()
	res := "[" + path.Base(file) + ":" + strconv.Itoa(line) + "]" + name
	tmp := strings.Split(name, ".")
	res = tmp[len(tmp)-1]
	return res
}

//var rand_gen = rand.New(rand.NewSource(time.Now().UnixNano()))
func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt() int {
	//return rand_gen.Int()
	return rand.Int()
}

func RandIntn(max int) int {
	//return rand_gen.Intn(max)
	return rand.Intn(max)
}

func NowInS() int64 {
	return time.Now().Unix()
}

func NowInNs() int64 {
	return time.Now().UnixNano()
}

func Abs(x int32) int32 {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0 // return correctly abs(-0)
	}
	return x
}

func Distance(flat float64, flng float64,
	tlat float64, tlng float64) (r int32) {
	distance := math.Sqrt((flat-tlat)*(flat-tlat) + (flng-tlng)*(flng-tlng))
	return int32(distance * 100000)
}
