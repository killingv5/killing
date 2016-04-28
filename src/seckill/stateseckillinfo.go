package seckill

import (
	"helpers/iowrapper"
	"time"
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

const (
	STATE_NOT_STARTED = 0
	STATE_ING         = 1
	STATE_ENDED       = 2
	STATE_NOT_EXIST   = 3
)

func (kp *Keeper) Run() {
	timediff := kp.Starttime.Sub(time.Now())
	time.Sleep(timediff)
	kp.State = STATE_ING
}

func ControlState(client *iowrapper.RedisClient) {
	for {
		time.Sleep(time.Second)
		infolist, err := GetAllProductInfo(client)
		if err != nil {
			logger.Error("GetAllProductInfo Failed! err=[%s]", err.Error())
			continue
		}
		for i := 0; i < len(infolist); i++ {
			pid := infolist[i].Pid
			_, ok := keepermap[pid]
			if ok {
				newstarttime := infolist[i].Seckillingtime
				t, _ :=time.Parse("20060102150405", newstarttime)
				// timediff := newstarttime.Sub(keepermap[pid].Starttime)
				timediff := t.Sub(keepermap[pid].Starttime)
				if timediff != 0 {
					delete(keepermap, pid)
					keepermap[pid] = &Keeper{STATE_NOT_STARTED, t}
					go keepermap[pid].Run()
				}
			} else {
				starttime := infolist[i].Seckillingtime
				t, _ :=time.Parse("20060102150405", starttime)
				keepermap[pid] = &Keeper{STATE_NOT_STARTED, t}
				go keepermap[pid].Run()
			}
		}
		for key, _ := range keepermap {
			infocount, err := GetProductCount(key, client)
			if err != nil {
				logger.Error("GetProductCount Failed! err=[%s]", err.Error())
				continue
			}
			if infocount <= 0 {
				keepermap[key].State = STATE_ENDED
			}
		}

	}
}
