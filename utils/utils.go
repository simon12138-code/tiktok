/*
* @Author: pzqu
* @Date:   2023/7/25 16:03
 */
package utils

import (
	"fmt"
	"time"
)

func GetNowFormatTodayTime() string {
	now := time.Now()
	dateStr := fmt.Sprintf("%02d-%02d-%02d", now.Year(), int(now.Month()), now.Day())
	return dateStr
}
