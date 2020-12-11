package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	co, _ := net.Dial("udp", "127.0.0.1:10001")
	go func() {
		var buf = make([]byte, 1024)
		for {
			n, _ := os.Stdin.Read(buf)
			co.Write(buf[:n])
		}
	}()
	var buf = make([]byte, 1024)
	for {
		n, _ := co.Read(buf)
		fmt.Println(string(buf[:n]))
	}

}
