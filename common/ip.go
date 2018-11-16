package common

import (
	"net"
	"os"
	"strings"
)

func GetIp() string {
	address, err := net.InterfaceAddrs()
	if err != nil {
		GetLog().Println("获取本机IP时异常退出", err)
		os.Exit(1)
	}
	var tempAds string
	for _, ads := range address {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := ads.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil && (strings.HasSuffix(ads.String(), "28") || strings.HasSuffix(ads.String(), "22")) { //子网掩码为22
				tempAds = ipNet.IP.String()
			}
		}
	}
	return tempAds
}
