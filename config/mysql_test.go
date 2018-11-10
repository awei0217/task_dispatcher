package config

import (
	"testing"
	"log"
	"fmt"
	"os"
	"bufio"
	"io"
	"strings"
)

func Benchmark_CreateConn(b *testing.B)  {

	for i:=0;i<10000;i++ {
		CreateConn();
		log.Println(i)
	}

}

func TestCreateConn(t *testing.T) {

	file ,_ :=os.Open("E:\\ccjf_price.sql")
	rd := bufio.NewReader(file)

	for {
		data, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		line := string(data)
		if strings.HasPrefix(line,"### DELETE FROM ccjf_price.pricestragy_price_definition"){
			sql := "INSERT INTO pricestragy_price_definition VALUES ("
			for {
				i := 0
				insert, _ := rd.ReadString('\n')
				insertLine := string(insert)
				if strings.HasPrefix(insertLine,"###   @") {
					if  strings.HasPrefix(insertLine,"###   @17") || strings.HasPrefix(insertLine,"###   @18"){
						sql = sql + "'"+strings.Split(insertLine,"=")[1]+"'" +","
					}else{
						sql = sql + strings.Split(insertLine,"=")[1] +","
					}

					i++
				}
				if strings.HasPrefix(insertLine,"###   @23") {
					sql = sql[0 : len(sql)-1] +");"
					break
				}

			}
			fmt.Println(sql)
		}

	}
}
