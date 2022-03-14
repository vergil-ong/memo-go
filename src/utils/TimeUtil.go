package utils

import (
	"math"
	"time"
)

func GetDiffDays(start, end int64) int {
	return int(math.Ceil((float64(end) - float64(start)) / float64(86400000)))
}

func GetTodayStartUnixMilli() int64 {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()).UnixMilli()
}

func GetTodayEndUnixMilli() int64 {
	oneDayMilli := 86400*1000 - 1
	return GetTodayStartUnixMilli() + int64(oneDayMilli)
}
