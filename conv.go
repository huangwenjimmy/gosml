package sml

import (
	"fmt"
	"strings"
	"time"
	"strconv"
	"bytes"
	"encoding/json"
	"reflect"
	"github.com/goinggo/mapstructure"
)

func ConvertDbVals(tp string, bs []byte) interface{}{
	if bs==nil{
		return nil
	}
	tp=strings.ToLower(tp)
	if strings.Contains(tp,"int"){
		val,_:= strconv.ParseInt(string(bs),10,64)
		return val
	}else if strings.EqualFold(tp,"double"){
		val,_:= strconv.ParseFloat(string(bs),64)
		return val
	}else if strings.Contains(tp,"date")||strings.Contains(tp,"time"){
		return ConvertStrToDate(string(bs))
	}else{
		return string(bs)
	}
}

func ConvertValueToRequestType(val string,tp string) interface{}{
	switch tp {
		case "char":
			return val
		case "number":
			if strings.Contains(val,"."){
				v,err:=strconv.ParseFloat(val,64)
				ThrowRuntime(err)
				return v
			}else{
				v,err:= strconv.ParseInt(val,10,64)
				ThrowRuntime(err)
				return v
			}
		case "array":
			return strings.Split(val,",")	
		case "date","time":
			return ConvertStrToDate(val)
		case "array-date","array-time":{
			vs:=strings.Split(val,",")
			times:=make([]time.Time,0)
			for _,v:=range vs{
				times=append(times,ConvertStrToDate(v))
			}
			return times
		}
		default:
			return val
					
	}
}

func ConvertStrToDate(val string) time.Time{
	bs:=[]byte(val)
	bf:=bytes.NewBufferString("");
	for _,v:=range bs{
		switch v{
			case '-',':',' ','/','T','Z':
				continue
			default:
				bf.WriteByte(v)
		}
	}
	format:=[]byte("200601021504050700")
	v,err:=time.Parse(string(format[:bf.Len()]),string(bf.Bytes()))
	ThrowRuntime(err)
	return v
}

func ThrowRuntime(err error){
	if err!=nil{
		panic (err)
	}
}

func ToJson(obj interface{}) ([]byte, error){
	return json.Marshal(obj)
}
func ToJsonString(obj interface{}) string{
	data,_:=ToJson(obj)
	return string(data)
}
func FromJson(data []byte,v interface{}) error{
	return json.Unmarshal(data, v)
}
func MapToStructure(m interface{}, rawVal interface{}) error{
	return mapstructure.Decode(m,rawVal)
}
func debug(a...interface{}){
	fmt.Println(a...)
}
func ConvertStructureToMap(pointer interface{}) map[string]interface{}{
	 result:=make(map[string]interface{})
	 ConvertStructureAppendMap(pointer,result)
     return result
}
func ConvertStructureAppendMap(pointer interface{},m map[string]interface{}){
	 elem := reflect.ValueOf(pointer).Elem()
	 relType := elem.Type()
     for i := 0; i < relType.NumField(); i++ {
        m[relType.Field(i).Name] = elem.Field(i).Interface()
     }
}
func ConvertToString(i interface{}) string{
	var result string
	if reflect.TypeOf(i)==nil{
		return result
	}
	kind:=reflect.TypeOf(i).Kind()
	rv:=reflect.ValueOf(i)
	switch{
		case kind==reflect.String: result=rv.String()
		case kind >= reflect.Int && kind <= reflect.Int64:result=strconv.FormatInt(rv.Int(),10)
		case kind >= reflect.Uint && kind <= reflect.Uint64: result=strconv.FormatUint(rv.Uint(),10)
		case kind >= reflect.Float32 && kind <= reflect.Float64:result=strconv.FormatFloat(rv.Float(), 'f', -1, 64)
		case kind == reflect.Bool: if rv.Bool(){result="1"}else{result="0"}
	}
	return result
}
func ConvertToInt(i interface{}) int64{
	var result int64
	if reflect.TypeOf(i)==nil{
		return result
	}
	kind:=reflect.TypeOf(i).Kind()
	rv:=reflect.ValueOf(i)
	switch{
		case kind==reflect.String: result,_=strconv.ParseInt(rv.String(),10,64)
		case kind >= reflect.Int && kind <= reflect.Int64:result=rv.Int()
		case kind >= reflect.Uint && kind <= reflect.Uint64: result=int64(rv.Uint())
		case kind >= reflect.Float32 && kind <= reflect.Float64:result=int64(rv.Float())
		case kind == reflect.Bool: if rv.Bool(){result=int64(1)}else{result=(int64)(0)}
	}
	return result
}
func ConvertToFloat(i interface{}) float64{
	var result float64
	if reflect.TypeOf(i)==nil{
		return result
	}
	kind:=reflect.TypeOf(i).Kind()
	rv:=reflect.ValueOf(i)
	debug(kind)
	switch{
		case kind==reflect.String: result,_=strconv.ParseFloat(rv.String(),64)
		case kind >= reflect.Int && kind <= reflect.Int64:result=float64(rv.Int())
		case kind >= reflect.Uint && kind <= reflect.Uint64: result=float64(rv.Uint())
		case kind >= reflect.Float32 && kind <= reflect.Float64:result=float64(rv.Float())
		case kind == reflect.Bool: if rv.Bool(){result=float64(1)}else{result=(float64)(0)}
	}
	return result
}

func ConvertStringToMap(evalStr string) map[string]string{
	m:=make(map[string]string)
	if !strings.HasPrefix(evalStr," "){
		evalStr=" "+evalStr
	}
	kvs:=strings.Split(evalStr," --")
	for _,v:=range kvs{
		if vs:=strings.SplitN(v,"=",2);len(vs)==2{
			m[vs[0]]=vs[1]
		}
	}
	return m
}
