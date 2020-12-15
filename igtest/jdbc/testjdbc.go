package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/huangwenjimmy/gosml"
)

func main() {
	stq := gosml.Init(gosml.Ds{Id: "defJt", Url: "root:root1234@tcp(192.168.1.234:3306)/xxxx", DbType: "mysql"})
	ds, _ := stq.Gdbcs["defJt"].QueryForList("select 1 as a")
	fmt.Println(ds)
}
