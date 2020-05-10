package gosml

import (
	"errors"
	//"bytes"
	"strings"
	"regexp"
	//"fmt"
	"reflect"
)

const(
	isNotEmpty="isNotEmpty"
	isEmpty="isEmpty"
	isNull="isNull"
	isNotNull="isNotNull"
)

var smlEl SmlEL
type Rst struct{
	DbType string 
	Sql string
	Args []interface{}
}

func (this *Rst) prettySql() string{
	return ""
}

func ParserSql(stp *SqlTemplate,gdbc *Gdbc) (*Rst,error){
	result:=&Rst{DbType:gdbc.DbType}
	args:=[]interface{}{}
	mainsql:=stp.Mainsql
	//
	re,_:=regexp.Compile("<("+isNotEmpty+"\\d*|"+isEmpty+"\\d*|"+isNull+"\\d*|"+isNotNull+"\\d*)\\s+property=\"(\\w+)\">")
	res:=re.FindAllString(mainsql,-1)
	for _,rev:=range res{//if
		start:=strings.Index(mainsql,rev)
		if start==-1{
			continue
		}
		mks:=re.FindStringSubmatch(rev)
		endMark:="</"+mks[1]+">"
		end:=strings.Index(mainsql,endMark)
		name:=mks[2]
		cberr:=checkConfigBuild(stp,name,mks[2])
		if cberr!=nil{
			return nil,cberr
		}
		content:=SubStr(mainsql,start+len(rev),end)
		currentAllContent:=SubStr(mainsql,start,end+len(endMark))
		flag:=checkValid(mks[1],name,stp)
		if flag{
			mainsql=strings.ReplaceAll(mainsql,currentAllContent,content)
		}else{
			mainsql=strings.ReplaceAll(mainsql,currentAllContent,"")
		}
	}
	re,_=regexp.Compile("<(if\\d*)\\s+test=\"([^\">]+)\">")
	res=re.FindAllString(mainsql,-1)
	for _,rev:=range res{//if
		start:=strings.Index(mainsql,rev)
		if start==-1{
			continue
		}
		mks:=re.FindStringSubmatch(rev)
		endMark:="</"+mks[1]+">"
		end:=strings.Index(mainsql,endMark)
		testEl:=mks[2]
		content:=SubStr(mainsql,start+len(rev),end)
		currentAllContent:=SubStr(mainsql,start,end+len(endMark))
		flag:=smlEl.Eval(testEl,stp.GetMap())
		if flag{
			mainsql=strings.ReplaceAll(mainsql,currentAllContent,content)
		}else{
			mainsql=strings.ReplaceAll(mainsql,currentAllContent,"")
		}
	}
	re,_=regexp.Compile("<(select\\d*|sql\\d*)\\s+id=\"(\\w+)\">");
	res=re.FindAllString(mainsql,-1)
	for _,rev:=range res{//if
		start:=strings.Index(mainsql,rev)
		mks:=re.FindStringSubmatch(rev)
		endMark:="</"+mks[1]+">"
		end:=strings.Index(mainsql,endMark)
		name:=mks[2]
		content:=SubStr(mainsql,start+len(rev),end)
		currentAllContent:=SubStr(mainsql,start,end+len(endMark))
		mainsql=strings.ReplaceAll(mainsql,currentAllContent,"")
		mainsql=strings.ReplaceAll(mainsql,"<ref id=\""+name+"\"/>",content)
	}
	re,_=regexp.Compile("\\$\\w+\\$");
	res=re.FindAllString(mainsql,-1)
	for _,rev:=range res{
		name:=SubStr(rev,1,len(rev)-1)
		value,err:=checkNotNull(name,stp)
		if err!=nil{
			return nil,err
		}
		mainsql=strings.ReplaceAll(mainsql,rev,ConvertToString(value))
	}
	re,_=regexp.Compile("#\\w+#");
	res=re.FindAllString(mainsql,-1)
	for _,rev:=range res{
		name:=SubStr(rev,1,len(rev)-1)
		value,err:=checkNotNull(name,stp)
		if err!=nil{
			return nil,err
		}
		mainsql=strings.ReplaceAll(mainsql,rev,"?")
		args=append(args,value)
	}
	re,_=regexp.Compile("\\s{2,}");
	mainsql=re.ReplaceAllString(mainsql," ")
	mainsql=strings.ReplaceAll(mainsql,"where 1=1 and","where")
	result.Args=args
	result.Sql=mainsql
	return result,nil
}

type item struct{
	start int
	startMark string
	end int
	endMark string
	content string
	flag bool
}

func checkConfigBuild(stp *SqlTemplate,name string,mark string) error{
	for _,sp:=range stp.SmlParams{
		if sp.Name==name{
			return nil;
		}
	}
	return errors.New(name +" is not config for "+mark)
}
func checkNotNull(name string,stp *SqlTemplate) (interface{},error){
	for _,sp:=range stp.SmlParams{
		if sp.Name==name&&IsNotEmpty(sp.Value){
			return sp.Value,nil;
		}
	}
	return nil,errors.New(name +" is config but is null ")
}


func checkValid(mark string,name string,stp *SqlTemplate) bool{
	for _,sp:=range stp.SmlParams{
		if sp.Name==name{
			val:=reflect.TypeOf(sp.Value)
			if(val==nil){
				return strings.HasPrefix(mark,"isEmpty")||strings.HasPrefix(mark,"isNull")
			}else{
				return (IsEmpty(sp.Value)&&strings.HasPrefix(mark, "isEmpty"))|| (IsNotEmpty(sp.Value)&&strings.HasPrefix(mark, "isNotEmpty"))
			}
		}
	}
	return false
}

