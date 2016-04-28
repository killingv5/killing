package seckill

import (
	"helpers/iowrapper"
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
	timediff := kp.starttime.Sub(time.Now())
	time.Sleep(timediff)
	kp.state = 1
}

func ControlState(client *iowrapper.RedisClient) {
	for {
		time.Sleep(time.Second)
		infolist, err := GetAllProductInfo(client)
		if err != nil {
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
					keepermap[pid] = &Keeper{0, newstarttime}
					go keepermap[pid].Run()
				}
			} else {
				starttime := infolist[i].seckillingtime
				keepermap[pid] = &Keeper{0, starttime}
				go keepermap[pid].Run()
			}
		}
		for key, _ := range keepermap {
			infocount, err := GetProductInfo(key, client)
			if err != nil {
				continue
			}
			if infocount <= 0 {
				keepermap[key].state = 2
			}
		}

	}
}
