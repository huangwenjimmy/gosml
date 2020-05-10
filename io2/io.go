package io2

import (
	"errors"
	"io"
	"bufio"
	"os"
)

var ERR_IGNORE=errors.New("ignore line")

type DoLine interface{
	ReadLine(row int,line string) error
}

func ReadFileLine(filename string,doLine DoLine) (int,error){
	file,err:=os.Open(filename)
	if err!=nil{
		return 0,err
	}
	return ReadLine(file,doLine)
}

func ReadLine(ir io.ReadCloser,doLine DoLine) (int,error){
	defer ir.Close()
	r:=bufio.NewReader(ir)
	var i int 
	for {
        line, err := readLine0(r)
        if err!=nil {
           if err==io.EOF{
	           break
           }else{
				return i,err	           
           }	
        }
        err2:=doLine.ReadLine(i, line)
        if(err2!=nil && err2!=ERR_IGNORE){
		      return i,err2
        }
        i++
    }
	return i,nil
}

func readLine0(r *bufio.Reader) (string, error) {
    line, isprefix, err := r.ReadLine()
    for isprefix && err == nil {
        var bs []byte
        bs, isprefix, err = r.ReadLine()
        line = append(line, bs...)
    }
    return string(line), err
}

