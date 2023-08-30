package cronjob

import "testing"

func TestReloadfiles(t *testing.T) {
	res, err := Reloadfiles("../public/", "record_video.txt")
	if err != nil {
		panic(err)
	}
	println(res)
}
