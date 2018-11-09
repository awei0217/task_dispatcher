package common

import (
	"encoding/json"
	"fmt"
	"testing"
)


func Test_SendMail(t *testing.T)  {

	m := make(map[string]string)
	m["id"]="37"
	m["tableName"]="trans_voucher201807"
	bs,_ := json.Marshal(m)
	fmt.Println(string(bs))
	SendMail(&Email{"sunpengwei1992@aliyun.com","golang测试邮件","朋伟：您好 \n \t这是golang的测试邮件"})

}


