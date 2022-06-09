package utils

import (
	"fmt"
	"log"
	"os"
)

func PrintLog(err error, errType string) {
	//在根目录创建一个logInfo.log文件用于存储日志信息
	logFile, ofErr := os.OpenFile("./logInfo.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 4)
	if ofErr != nil {
		fmt.Println("open logfile failed, err:", err)
		return
	}
	//将日志输出到logInfo.log文件
	log.SetOutput(logFile)
	//设置日志的前缀信息,即错误的类型：Error，fatal，普通输出info等
	log.SetPrefix(errType)
	//设置日志的输出格式
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.Println(err)
}
