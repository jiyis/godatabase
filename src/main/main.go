package main

import (
	. "connection"
	. "library"
)

func main() {
	//初始化日志类
	InitLog()
	//初始化数据库链接
	Init()
	//处理liquibase
	ProcessLiquibase()
	
}
