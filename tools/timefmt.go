package tools

import (
	"fmt"
	"time"
)

// Now2mysql
// Convert time.now() to mysql datetime format
func Now2mysql() (mysqlDatetime string) {
	now := time.Now()
	mysqlDatetimeFormat := "2006-01-02 15:04:05"
	mysqlDatetime = now.Format(mysqlDatetimeFormat)
	return
}

func Time2mysql(t time.Time) (mysqlDatetime string) {
	mysqlDatetimeFormat := "2006-01-02 15:04:05"
	mysqlDatetime = t.Format(mysqlDatetimeFormat)
	return
}

func Mysql2Unix(t string) (unixTimestamp int64) {
	timeObj, err := time.Parse(time.RFC3339, t)
	if err != nil {
		fmt.Println(err)
	}
	unixTimestamp = timeObj.Unix()
	return
}
