package sml

import (
	"database/sql"
	//"fmt"
	"errors"
	"log"
)

var gdbcLog bool=true

type Gdbc struct{
	Db *sql.DB
	DbType string
	Convert func(ty string,val []byte) interface{}
}


func(this *Gdbc) QueryForMap(sql string,args...interface{}) (map[string]interface{},error){
	i:=0;
	var result map[string]interface{}
	err:=this.QueryForCall(sql,func(m map[string]interface{}) (bool,error){
			i++;
			result=m;
			if i==2{
				return true,errors.New("more rows error")
			}
			return true,nil
	},args...)
	if err!=nil{
		return nil,err
	}
	if i==0{
		return nil,errors.New("not exists rows")
	}
	return result,nil
}


func(this *Gdbc) QueryForList(sql string,args ... interface{}) ([]map[string]interface{},error){
	result:=make([]map[string]interface{},0)
	err:=this.QueryForCall(sql,func(m map[string]interface{}) (bool,error){
		result=append(result,m);
		return true,nil
	},args...)
	if err!=nil{
		return nil,err;
	}
	return result,nil;
}

func (this *Gdbc) QueryForCall(sql string,hand func(map[string]interface{}) (bool,error),args... interface{}) (error){
	if gdbcLog{
		log.Println("sql:[%s],%v",sql,args)
	}
	rows,err := this.Db.Query(sql,args...)
	if err!=nil{
		return err
	}
	cns,_:=rows.Columns()
    vals := make([][]byte, len(cns))
    scans :=make([]interface{},len(cns))
    for k,_:=range vals{
	    scans[k]=&vals[k]
    }
    vts,_:=rows.ColumnTypes()
    for rows.Next(){
	    rows.Scan(scans...);
	    r:=make(map[string]interface{})
	    for i:=0;i<len(vals);i++{
		    r[cns[i]]=this.Convert(vts[i].DatabaseTypeName(),vals[i])
	    }
	    ok,err:=hand(r)
	    if !ok{
		    break;
	    }
	    if err!=nil{
		    return err
	    }
    }
	return nil
}


func EnabledLog(enabled bool){
	gdbcLog=enabled
}
