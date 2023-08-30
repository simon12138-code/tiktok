package asyncJob

import (
	"go_gin/global"
	"testing"
	"time"
)

func TestColletorVideoinfo(t *testing.T) {
	global.VideoChan = make(chan global.VideoInfo, 10)
	stop := make(chan struct{})
	//协程测试一定要保证异步
	go func() {
		time.Sleep(time.Second * 5)
		for i := 0; i < 101; i++ {
			global.VideoChan <- global.VideoInfo{VideoId: 1, VideoFileName: "test1"}
		}
	}()
	ColletorVideoInfo(stop)
}
