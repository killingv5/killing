package seckill

import (
	"fmt"
	"time"
)

var keepermap map[string]*Keeper

func init() {
	keepermap = map[string]*Keeper{}
}

type Keeper struct {
	state     int
	starttime time.Time
}

func (kp *Keeper) Run() {
	timediff := -time.Now().Sub(kp.starttime)
	fmt.Printf("time:%v\n", timediff)
	time.Sleep(timediff)
	kp.state = 1
}

//func StartKeeper(pid string) {
//	for;; {
//		time.Sleep(0.01 * time.Second)
//		info, err := GetProductInfo(pid)
//		timediff := time.Equal(info.seckillingtime, time.Now())
//		if timediff {
//			keepermap[pid] = 1
//			break
//		}
//	}
//}

//func EndKeeper(pid string) {
//	for;; {
//		time.Sleep(0.02 * time.Second)
//		info, err := GetProductInfo(pid)
//		if info.Pnum <= 0 {
//			keepermap[pid] = 2
//			break
//		}
//	}
//}
type productTest struct {
	Pid            string
	seckillingtime time.Time
	Pnum           int64
}

func ControlState1() {
	infolist := []productTest{
		productTest{
			Pid:            "1111",
			seckillingtime: time.Now().Add(time.Second * 3),
			Pnum:           int64(100)},
		productTest{
			Pid:            "2222",
			seckillingtime: time.Now().Add(time.Second * time.Duration(45) / 10),
			Pnum:           int64(100)}}
	sign := 0
	for {
		time.Sleep(time.Second)
		//infolist, err := GetAllProductInfo()
		//if err != nil {
		//	continue
		//}
		for i := 0; i < len(infolist); i++ {
			pid := infolist[i].Pid
			_, ok := keepermap[pid]
			if ok {
				newstarttime := infolist[i].seckillingtime
				timediff := newstarttime.Sub(keepermap[pid].starttime)
				if timediff != 0 {
					delete(keepermap, pid)
					keepermap[pid] = &Keeper{0, newstarttime}
					go keepermap[pid].Run()
				}
			} else {
				starttime := infolist[i].seckillingtime
				keepermap[pid] = &Keeper{0, starttime}
				go keepermap[pid].Run()
			}
		}
		//for key, value := range keepermap {
		//	info, err := GetProductInfo(key)
		//	if err != nil {
		//		continue
		//	}
		//	if info.Pnum <= 0 {
		//		keepermap[key].state = 2
		//		fmt.Printf("end:%v,2\n", key)
		//	}
		//}
		for i := 0; i < len(infolist); i++ {
			pid := infolist[i].Pid
			if infolist[i].Pnum <= 0 {
				keepermap[pid].state = 2
				fmt.Printf("end:%v, 2\n", pid)
			}
		}
		for i := 0; i < len(infolist); i++ {
			pid := infolist[i].Pid
			fmt.Printf("state:%v, %v, %+v, %+v, %v\n", time.Now(), pid, keepermap[pid].state, keepermap[pid].starttime, infolist[i].Pnum)
		}
		for key, _ := range keepermap {
			fmt.Printf("mapstate:%v, %v, %+v, %+v\n", time.Now(), key, keepermap[key].state, keepermap[key].starttime)
		}

		if len(infolist) < 3 {
			infolist = append(infolist, productTest{
				Pid:            "3333",
				seckillingtime: time.Now().Add(time.Second * 5),
				Pnum:           int64(100)})
		}
		if sign == 0 {
			infolist[0].seckillingtime = time.Now().Add(time.Second * 7)
			sign = 1
		}

	}
}
