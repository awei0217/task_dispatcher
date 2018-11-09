package config

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose"
	"task_dispatcher/common"
)


var DbConfig = map[string]interface{}{
	// Default database configuration
	"Default": "mysql_dev",
	// (Connection pool) Max open connections, default value 0 means unlimit.
	"SetMaxOpenConns": 300,
	// (Connection pool) Max idle connections, default value is 1.
	"SetMaxIdleConns": 10,

	// Define the database configuration character "mysql_dev".
	"Connections": map[string]map[string]string{
		"mysql_dev": map[string]string{
			"host":     "127.0.0.1",
			"username": "root",
			"password": "root",
			"port":     "3306",
			"database": "scheduler",
			"charset":  "utf8",
			"protocol": "tcp",
			"prefix":   "",      // Table prefix
			"driver":   "mysql", // Database driver(mysql,sqlite,postgres,oracle,mssql)
		},
	},
}

var Conn gorose.Connection
func init()  {
	connection , err := gorose.Open(DbConfig)
	if err != nil {
		common.GetLog().Errorln("初始化数据库报错",err)
		panic(err)
	}
	Conn = *connection
}

func CreateConn() (gorose.Connection)  {
	connection , _ := gorose.Open(DbConfig)
	return *connection
}
