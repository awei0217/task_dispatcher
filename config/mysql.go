package config

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
	"task_dispatcher/common"
)

var dbConfig = &gorose.DbConfigSingle{
	Driver:          "mysql",                                                // driver: mysql/sqlite/oracle/mssql/postgres
	EnableQueryLog:  true,                                                   // if enable sql logs
	SetMaxOpenConns: 5,                                                      // connection pool of max Open connections, default zero
	SetMaxIdleConns: 1,                                                      // connection pool of max sleep connections
	Prefix:          "",                                                     // prefix of table
	Dsn:             "root:root@tcp(127.0.0.1:3306)/scheduler?charset=utf8", // db dsn
}
var Conn gorose.Connection

func init() {
	connection, err := gorose.Open(dbConfig)
	if err != nil {
		common.GetLog().Errorln("初始化数据库报错", err)
		panic(err)
	}
	Conn = *connection
}

func CreateConn() gorose.Connection {
	connection, _ := gorose.Open(dbConfig)
	return *connection
}
