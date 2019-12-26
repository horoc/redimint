package core

import (
	"github.com/chenzhou9513/redimint/database"
	"github.com/jasonlvhit/gocron"
)

var Schedulers *gocron.Scheduler

func InitAllJobs() {
	Schedulers = gocron.NewScheduler()
	Schedulers.Every(10).Seconds().Do(CheckRedisStatus)
}

func CheckRedisStatus() {
	isAlive := database.CheckAlive(3)
	if !isAlive {
		err := AppService.RestoreLocalDatabase()
		if err != nil {
			return
		}
	}
}

func StartAllJobs(){
	Schedulers.Start()
}

func StopAllJobs(){
	Schedulers.Clear()
}
