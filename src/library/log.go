package library

import (
	"os"
	"log"
)

var Logger *log.Logger

func InitLog() {
	fileName := "error.log"
	logFile, err := os.Create(fileName)
	
	if err != nil {
		log.Fatalln("log open file error !")
	}
	Logger = log.New(logFile, "[Debug]", log.Llongfile)
	
	//debugLog.Println("A debug message here")
	//debugLog.SetPrefix("[Info]")
	//debugLog.Println("A Info Message here ")
	//debugLog.SetFlags(debugLog.Flags() | log.LstdFlags)
	//debugLog.Println("A different prefix")
}
