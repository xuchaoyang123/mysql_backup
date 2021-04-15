package logfile


import (
	"log"
	"os"
	"time"
)

var (
	Loger *log.Logger
)

func Print2() {

	//日志生成Loger
	//mysql := config.ConfMap["mysql"]

	FileName := "./check_list_" + time.Now().Format("2006-01-02") + ".log"

	logFile, err := os.OpenFile(FileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic("这")
	}

	Loger = log.New(logFile, "", log.LstdFlags) // 将文件设置为loger作为输出

	return

}

