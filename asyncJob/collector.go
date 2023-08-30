package asyncJob

import (
	"bufio"
	"fmt"
	"go_gin/global"
	"os"
	"path"
	"time"
)

// 异步存储Video信息方便后期重整Videoinfo
func ColletorVideoInfo(stopchan chan struct{}) {
	videoInfo := make([]global.VideoInfo, 0, global.VideoInfoCollectorMaxNum)
	for {
		select {
		//阻塞性接受通道的信息，感觉该方案也可以使用消息队列
		case v := <-global.VideoChan:
			//容量到达限额，进行转储
			if len(videoInfo) == global.VideoInfoCollectorMaxNum {
				//文件操作
				for i := 0; i < global.CollectorRetryTime; i++ {
					err := SaveInFile(videoInfo, "record_video.txt")
					if err != nil {
						global.Lg.Error(err.Error())
						time.Sleep(global.CollectorRetryTimeDuration)
						//失败重新尝试
						continue
					} else {
						//成功退出循环
						//重置缓存
						videoInfo = make([]global.VideoInfo, 0, global.VideoInfoCollectorMaxNum)
						break
					}
				}
			}
			//缓存
			videoInfo = append(videoInfo, global.VideoInfo{VideoId: v.VideoId, VideoFileName: v.VideoFileName})
		case <-stopchan:
			break
		default:
		}

	}
}
func ColletorCoverInfo(stopchan chan struct{}) {
	coverInfo := make([]global.VideoInfo, 0, global.VideoInfoCollectorMaxNum)
	for {
		select {
		//阻塞性接受通道的信息，感觉该方案也可以使用消息队列
		case v := <-global.CoverChan:
			//容量到达限额，进行转储
			if len(coverInfo) == global.VideoInfoCollectorMaxNum {
				//文件操作
				for i := 0; i < global.CollectorRetryTime; i++ {
					err := SaveInFile(coverInfo, "record_cover.txt")
					if err != nil {
						global.Lg.Error(err.Error())
						time.Sleep(global.CollectorRetryTimeDuration)
						//失败重新尝试
						continue
					} else {
						//成功退出循环
						//重置缓存
						coverInfo = make([]global.VideoInfo, 0, global.VideoInfoCollectorMaxNum)
						break
					}
				}
			}
			//缓存
			coverInfo = append(coverInfo, global.VideoInfo{VideoId: v.VideoId, VideoFileName: v.VideoFileName})
		case <-stopchan:
			break
		default:
		}

	}
}

func SaveInFile(infos []global.VideoInfo, fileName string) error {
	file, err := os.OpenFile(path.Join("../public/", fileName), os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	//关闭文件
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, v := range infos {
		str := fmt.Sprintf("%d_____%s\n", v.VideoId, v.VideoFileName)
		println(str)
		_, err = writer.WriteString(str)
		if err != nil {
			return err
		}
	}
	//从缓冲刷到文件中
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
