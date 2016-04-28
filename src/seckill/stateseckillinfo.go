package seckill

import (
	"helpers/iowrapper"
	"time"
)

var keepermap map[string]*Keeper

func init() {
	keepermap = map[string]*Keeper{}
}

func GetPidState(pid int) int {
	if val, ok := keepermap[fmt.Printf("%v", pid)]; ok {
		return val.state
	}
	return STATE_NOT_EXIST
}

const (
	STATE_NOT_STARTED = 0
	STATE_ING         = 1
	STATE_ENDED       = 2
	STATE_NOT_EXIST   = 3
)

type Keeper struct {
	state     int
	starttime time.Time
}

func (kp *Keeper) Run() {
	timediff := kp.starttime.Sub(time.Now())
	time.Sleep(timediff)
	kp.state = STATE_ING
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
				newstarttime := infolist[i].seckillingtime
				timediff = newstarttime.Sub(keepermap[pid].starttime)
				if timediff != 0 {
					delete(keepermap, pid)
					keepermap[pid] = &Keeper{STATE_NOT_STARTED, newstarttime}
					go keepermap[pid].Run()
				}
			} else {
				starttime := infolist[i].seckillingtime
				keepermap[pid] = &Keeper{STATE_NOT_STARTED, starttime}
				go keepermap[pid].Run()
			}
		}
		for key, _ := range keepermap {
			infocount, err := GetProductInfo(key, client)
			if err != nil {
				logger.Error("GetProductInfo Failed! err=[%s]", err.Error())
				continue
			}
			if infocount <= 0 {
				keepermap[key].state = STATE_ENDED
			}
		}

	}
}
