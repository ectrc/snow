package aid

import "time"

func TimeStartOfDay() string {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location()).Format("2006-01-02T15:04:05.999Z")
}

func TimeEndOfDay() string {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 999999999, time.Now().Location()).Format("2006-01-02T15:04:05.999Z")
}

func TimeEndOfWeekString() string {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 999999999, time.Now().Location()).AddDate(0, 0, 7).Format("2006-01-02T15:04:05.999Z")
}