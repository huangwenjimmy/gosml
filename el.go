package sml

import (
	"strings"
	"regexp"
)

type EL interface{
	Eval(evalStr string,context map[string]interface{}) interface{}
}

type SmlEL struct{
	
}

func (this *SmlEL) Eval(evalStr string,context map[string]interface{}) bool{
	m:=ConvertStringToMap(evalStr)
	flag:=true
	isOr:=strings.EqualFold(m["this.type"],"or")
	for k,v:=range m{
		if(k=="this.type"){
			continue
		}
		if vs:=strings.SplitN(v," ",2);len(vs)==2{
			flag=OperatorEval(vs[0],vs[1],ConvertToString(context[k]))
			if(isOr&&flag){
				return true
			}
			if(!isOr && !flag){
				return false
			}
		}
	}
	return flag
}

func OperatorEval(operator string, vc string, realV string) bool{
	result:=false
	if len(realV)==0&&operator!="is"{
		return result
	}
	switch operator{
		case "=","eq":result=vc==realV
		case "!=","ne":result=vc!=realV
		case "is": if(vc=="nil"){result=len(realV)==0}else if(vc=="notnil"){result=len(realV)>0}
		case ">","gt" :result=ConvertToFloat(realV)>ConvertToFloat(vc)
		case ">=","ge" :result=ConvertToFloat(realV)>=ConvertToFloat(vc)
		case "<","lt":result=ConvertToFloat(realV)<ConvertToFloat(vc)
		case "<=","le":result=ConvertToFloat(realV)<=ConvertToFloat(vc)
		case "like","contain":result=strings.Contains(realV,vc)
		case "ilike":result=strings.Contains(strings.ToLower(realV),strings.ToLower(vc))
		case "nlike":result=!strings.Contains(realV,vc)
		case "nilike":result=!strings.Contains(strings.ToLower(realV),strings.ToLower(vc))
		case "in":result=Contains(strings.Split(vc,","),realV)
		case "regexp":if rd,err:=regexp.MatchString(vc,realV);err==nil{result=rd;debug(rd)}
		case "nregexp":if rd,err:=regexp.MatchString(vc,realV);err==nil{result=!rd}
	}
	return result
}