package net_server

import (
	"net"
	"task_dispatcher/common"
	"bufio"
)
var tcpListener *net.TCPListener
func InitServer()  {
	go StartServer()
}
func StartServer()  {
	tcpAddr,_ := net.ResolveTCPAddr("tcp4",common.GetIp()+":10010")
	tcpListener,_ := net.ListenTCP("tcp",tcpAddr)
	common.GetLog().Println(common.GetIp(),"socket服务端启动成功")
	for{
		con,err := tcpListener.Accept();
		if err != nil {
			continue
		}
		go readNetInfo(con)
	}
}
//TODO 待完善
func StopServer()  {
	common.GetLog().Println(common.GetIp(),"socket服务端停止成功")
}

func readNetInfo(conn net.Conn)  {
	r :=bufio.NewReader(conn)
	str,_:=r.ReadString('@')
	common.GetLog().Infoln("收到IP:",conn.RemoteAddr().String(),"的心跳,内容为:",str)
	HandlerNetMessage(string(str[0:len(str)-1]),conn)
	defer func() {
		conn.Close()
	}()
}
