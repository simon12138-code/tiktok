package initialize

import (
	"go_gin/asyncJob"
	"go_gin/global"
)

func InitColloctor() {
	global.VideoChan = make(chan global.VideoInfo, 10)
	s := make(chan struct{})
	asyncJob.ColletorVideoInfo(s)
	asyncJob.ColletorCoverInfo(s)
}
