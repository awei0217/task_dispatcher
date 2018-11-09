package net_client

import (

	"testing"
)
func TestSendNetInfo(t *testing.T) {
	SendNetInfo("127.0.0.1","您好,服务端")
}