package net_client

import (
	"net"
	"task_dispatcher/common"
	"bufio"
)
func SendNetInfo(ip string,message string) (int,error){

	tcpAddr,_ := net.ResolveTCPAddr("tcp4",ip+":10010")
	con,err1 :=net.DialTCP("tcp",nil,tcpAddr)
	if err1 != nil {
		common.GetLog().Errorln("连接服务端出错",err1)
		return 0,err1
	}
	w := bufio.NewWriter(con)
	count,err2 := w.WriteString(message+"@")
	w.Flush()
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	return count,err2

}