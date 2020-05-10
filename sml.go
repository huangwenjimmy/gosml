package gosml

import (
	"log"
	"database/sql"
	"fmt"
)


type Sml struct{
	Dss map[string]*sql.DB
	Gdbcs map[string]*Gdbc
	Cm CacheManager
}

type Ds struct{
	Id string
	Url string
	DbType string
	MaxConn int
	MaxIdle int
}

func (this *Sml) Cache(cm CacheManager) *Sml{
	this.Cm=cm
	return this;
}

func  Init(dss...Ds) *Sml{
	dvv:= map[string]*sql.DB{}
	gds:= map[string]*Gdbc{}
	sml:=&Sml{Dss:dvv,Gdbcs:gds}
	for _,ds:=range dss{
		db, err := sql.Open(ds.DbType,ds.Url)
		if err!=nil{
			log.Fatalf("init %s url:%s error:%s",ds.Id,ds.Url,err)	
			continue
		}
		if ds.MaxConn>0{
			db.SetMaxIdleConns(ds.MaxConn)
		}
		if ds.MaxIdle>0{
			db.SetMaxOpenConns(ds.MaxIdle)
		}
		sml.Dss[ds.Id]=db
		sml.Gdbcs[ds.Id]=&Gdbc{Db:db,Convert:ConvertDbVals,DbType:ds.DbType}
	}
	return sml
}

func (this *Sml) GetSqlTemplate(ifId string) (*SqlTemplate,error){
	key:=fmt.Sprintf("jdbc:%s:sqltemplate",ifId)
	bst:=this.Cm.Get(key)
	st:=&SqlTemplate{DbId:"defJt"}
	if len(bst)>0{
		unerr:=FromJson(bst,st)
		return st,unerr
	}
	stmp,err:=this.Gdbcs["defJt"].QueryForMap("select id Id,mainsql Mainsql,condition_info as ConditionInfo,rebuild_info as RebuildInfo,cache_min CacheMin,descr,update_time UpdateTime,db_id DbId from sml_if where id=?",[]interface{}{ifId}...)
	if err!=nil{
		return nil,err
	}
	MapToStructure(stmp,st)
	st.Init()
	ss,err:=ToJson(st)
	this.Cm.Put(key,ss,-1)
	return st,nil
}

func (this *Sml) Query(ifId string,params map[string]string) ([]map[string]interface{},error){
	stp,_:=this.GetSqlTemplate(ifId)
	stp.ReinitParam(params)
	rst,err:=ParserSql(stp,this.Gdbcs[stp.DbId])
	if err!=nil{
		return nil,err
	}
	return this.Gdbcs[stp.DbId].QueryForList(rst.Sql,rst.Args...)
}


//func (sml *Sml) Update()
