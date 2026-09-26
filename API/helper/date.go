package helper

import "time"

func CurrentDate() string {
	return time.Now().Format("2006-01-02")
}

func CurrentTime() string {
	loc, _ := time.LoadLocation("Asia/Phnom_Penh")
	return time.Now().In(loc).Format("15:04:05")
}
