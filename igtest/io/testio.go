package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gosml"
	"gosml/ftps"
	"gosml/io2"
	"io/ioutil"
	"os"
	"os/exec"
)

type Dod struct{}

func (this *Dod) ReadLine(row int, line string) error {
	fmt.Println(row, "---", line)
	return nil
}

func main1() {
	io2.ReadFileLine("d:/temp/json/map.txt", &Dod{})
}
func main2() {
	url := ftps.NewFtps("ftp://192.168.1.234:21/a/test.txt?timeout=10")
	fmt.Println(url)
	ftps := ftps.NewFtps("ftp://smp:smp1234@192.168.1.234:21//home/smp/upload/test1.txt?timeout=1000").Login()
	fmt.Println(ftps.Conn.CurrentDir())
	f, _ := os.Open("d:/temp/json/map.txt")
	//ftps.Conn.ChangeDir("upload")
	//ftps.Put(f)
	ir, _ := ftps.Get()
	bs, _ := ioutil.ReadAll(ir)
	fmt.Println(string(bs))
	defer f.Close()
	defer ftps.DisConnect()
}
func main() {
	cmd := exec.Command("ping", "wwww.baidu.com")
	ic, err := cmd.Output()
	gosml.ThrowRuntime(err)
	a := transform.NewReader(bytes.NewReader(ic), simplifiedchinese.GBK.NewEncoder())
	bs, _ := ioutil.ReadAll(a)
	fmt.Println(string(bs))
}
