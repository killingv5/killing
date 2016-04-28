package seckill

import (
	"helpers/iowrapper"
	"time"
	// "fmt"

	logger "github.com/xlog4go"
)

type Keeper struct {
	State     int
	Starttime time.Time
}

var keepermap map[string]*Keeper

func init() {
	keepermap = map[string]*Keeper{}
}

func GetPidState(pid string) int {
	if val, ok := keepermap[pid]; ok {
		return val.State
	}
	return STATE_NOT_EXIST
}

func (kp *Keeper) Run() {
	timediff := kp.Starttime.Sub(time.Now())
	// fmt.Printf("state update time:[%+v]\n",time.Now())
	// fmt.Printf("state update time:[%+v]\n",kp.Starttime)
	// fmt.Printf("state update after:[%+v]\n",timediff)
	time.Sleep(timediff)
	// fmt.Printf("state update count:[%+v]",kp)

	kp.State = STATE_ING
}

func ControlState(client *iowrapper.RedisClient) {
	for {
		time.Sleep(time.Second)
		// fmt.Println("control count")
		infolist, err := GetAllProductInfo(client)
		if err != nil {
			logger.Error("GetAllProductInfo Failed! err=[%s]", err.Error())
			continue
		}
		infomap := make(map[string]int)
		for i := 0; i < len(infolist); i++ {
			pid := infolist[i].Pid
			infomap[pid] = 1
			_, ok := keepermap[pid]
			if ok {
				newstarttime := infolist[i].Seckillingtime
				//t, _ := time.Parse("20060102150405", newstarttime)
				timediff := newstarttime.Sub(keepermap[pid].Starttime)
				//timediff := t.Sub(keepermap[pid].Starttime)
				if timediff != 0 {
					delete(keepermap, pid)
					//keepermap[pid] = &Keeper{STATE_NOT_STARTED, t}
					keepermap[pid] = &Keeper{STATE_NOT_STARTED, newstarttime}
					go keepermap[pid].Run()
				}
			} else {
				starttime := infolist[i].Seckillingtime
				//t, _ := time.Parse("20060102150405", starttime)
				//keepermap[pid] = &Keeper{STATE_NOT_STARTED, t}
				keepermap[pid] = &Keeper{STATE_NOT_STARTED, starttime}
				go keepermap[pid].Run()
			}
		}
		for key, _ := range keepermap {
			_, ok := infomap[key]
			if !ok {
				delete(keepermap, key)
				continue
			}
			//infocount, err := GetProductCount(key, client)
			//if err != nil {
			//	logger.Error("GetProductCount Failed! err=[%s]", err.Error())
			//	continue
			//}
			//if infocount <= 0 {
			//	keepermap[key].State = STATE_ENDED
			//}
		}

	}
}
