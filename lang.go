package gosml

import (
	"strings"
	"reflect"
	"os"
)

type Pairs map[string]string



func init(){
	getMapArgs0()
}


func getMapArgs0() Pairs{
	result:=make(map[string]string)
	for _,arg:=range os.Args{
		if strings.HasPrefix(arg,"--")&&strings.Contains(arg,"="){
			kv:=strings.SplitN(arg,"=",2)
			result[kv[0]]=kv[1]
		}
	}
	return result
}


func SubStr(source string, start int, end int) string {
    var r = []rune(source)
    length := len(r)
    if start < 0 || end > length || start > end {
        return ""
    }
    if start == 0 && end == length {
        return source
    }
    return string(r[start : end])
}


func ContainKeyCaseInsensitive(m map[string]interface{}, key string) bool{
	for k,_:=range m{
		if strings.EqualFold(k,key){
			return true
		}
	}
	return false
}

func ContainKey(m map[string]interface{},key string) bool{
	for k,_:=range m{
		if k==key{
			return true
		}
	}
	return false
}

func IsEmpty(i interface{}) bool{
	if reflect.TypeOf(i)==nil{
		return true
	}
	rv:=reflect.ValueOf(i)
	kind:=rv.Kind()
	result:=false
	debug(kind)
	switch{
		case kind==reflect.String: result=len(rv.String())==0
		case kind==reflect.Slice : result=rv.Len()==0
		case kind==reflect.Array: result=rv.Len()==0
	}
	return result
}

func IsNotEmpty(i interface{}) bool{
	return !IsEmpty(i)
}

func Contains(array []string,ele string) bool{
	for _,v:=range array{
		if(v==ele){
			return true
		}
	}
	return false
}

