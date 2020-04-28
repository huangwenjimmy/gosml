package sml

import (
	"errors"
	"io"
)

var ERR_IGNORE=errors.New("ignore line")

type DoLine interface{
	ReadLine(row int,line string) error
}

func ReadLine(ir io.ReadCloser,charset string,doLine DoLine){
	
}

