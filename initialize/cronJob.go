package initialize

import (
	"go_gin/cronjob"
	"go_gin/global"
	"time"
)

func InitCronJob() {
	stop := make(chan struct{})
	cron := cronjob.NewCronJobList(stop)
	timer := time.NewTicker(global.MaxUrlExpireTime)
	job := cronjob.Job{Ticker: *timer, Func: cronjob.HttpReset}
	err := cron.RigisterJob("updateUrl", job)
	if err != nil {
		panic(err)
	}
	cron.Run()
}
