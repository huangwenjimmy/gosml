package main

import (
	"fmt"
	"github.com/huangwenjimmy/gosml"
	"github.com/huangwenjimmy/gosml/net2"
	"net"
	"strings"
)

func main1() {
	ln, err := net.Listen("tcp", "127.0.0.1:10001")
	gosml.ThrowRuntime(err)
	for {
		conn, err2 := ln.Accept()
		fmt.Println(conn.RemoteAddr())
		gosml.ThrowRuntime(err2)
		go handleConn(conn)
	}
}
func main2() {
	ud, _ := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	handleUdpConn(ud)
}
func main() { //iprange
	fmt.Println(net2.IpToInt64("192.168.1.120"))
	fmt.Println(net2.Int64ToIp(3232235896))
	vs := net2.NewIpv4Ranges("192.54.1.0/24,192.54.1.0/16")
	fmt.Println(vs.Contains("192.54.11.123"))
	//fmt.Println(net2.Int64ToIp(vs.Start))
	//fmt.Println(net2.Int64ToIp(vs.End))
}

func handleUdpConn(addr *net.UDPAddr) {
	ul, _ := net.ListenUDP("udp", addr)
	var buf = make([]byte, 1024)
	for {
		n, radd, _ := ul.ReadFromUDP(buf)
		if string(buf[:n]) == "byte" {
			break
		}
		ul.WriteToUDP([]byte(strings.ToUpper(string(buf[:n]))), radd)
	}
}
func handleConn(conn net.Conn) {
	dst, err := net.Dial("tcp", "192.168.1.120:10001")
	gosml.ThrowRuntime(err)
	go swap(dst, conn)
	swap(conn, dst)
}
func swap(r, w net.Conn) {
	defer r.Close()
	defer w.Close()
	var buf = make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil {
			break
		}
		_, err = w.Write(buf[0:n])
		if err != nil {
			break
		}
	}
}
