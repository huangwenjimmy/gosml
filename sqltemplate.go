package gosml

import (
	"strings"
	"time"
)

type SqlTemplate struct{
	Id string 
	Mainsql string
	ConditionInfo string
	RebuildInfo string
	CacheMin int
	Descr string
	DbId string
	UpdateTime time.Time
	SmlParams []*SmlParam
	Rebuild *Build
}


func (this *SqlTemplate) Init(){
	this.SmlParams=make([]*SmlParam,0)
	if len(this.ConditionInfo)>0{
		sepe:=","
		for i,v:=range strings.Split(this.ConditionInfo,"\n"){
			v=strings.TrimSpace(v)
			if len(v)==0{
				continue
			}
			if i==0&&strings.HasPrefix(v,"#"){
				sepe=strings.ReplaceAll(v,"#","")
				continue
			}
			if strings.HasPrefix(v,"#"){
				continue
			}
			smlParam:=&SmlParam{Type:"char"}
			for j,vt:=range strings.Split(v,sepe){
				switch j {
					case 0:
						smlParam.Name=vt
					case 1:
						smlParam.Type=vt
					case 2:
						smlParam.DefaultValue=vt
					case 3:
						smlParam.Type=vt	
				}
			}
			this.SmlParams=append(this.SmlParams,smlParam)
		}
	}
	this.Rebuild=&Build{Name:"list",Options:make(map[string]string)}
	
	if len(this.RebuildInfo)>0{
		for i,v:=range strings.Split(this.RebuildInfo,"\n"){
			if i==0{
				this.Rebuild.Name=v
			}else{
				kv:=strings.SplitN(v,"=",2)
				if len(kv)==2{
					this.Rebuild.Options[kv[0]]=kv[1]
				}
				
			}
		}
	}
}

func (this *SqlTemplate) ReinitParam(params map[string]string){
	sps:=this.SmlParams
	for _,v:=range sps{
		val:=params[v.Name]
		if len(val)>0{
			v.Value=ConvertValueToRequestType(val,v.Type)
		}else if(len(v.DefaultValue)>0){
			v.Value=ConvertValueToRequestType(v.DefaultValue,v.Type)
		}
	}
}

func (this *SqlTemplate) GetMap() map[string]interface{}{
	result:=make(map[string]interface{})
	for _,v:=range this.SmlParams{
		result[v.Name]=v.Value
	}
	return result
}


type SmlParam struct{
	Name string
	Value interface{}
	Type string
	DefaultValue string
	Descr string
}

type Build struct{
	Name string
	Options map[string] string
}