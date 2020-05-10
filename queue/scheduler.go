package queue

import (
	"time"
	"strings"
	"regexp"
	"strconv"
	"errors"
	"github.com/huangwenjimmy/gosml"
	"sync"
)
var months= [12]string{"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"}
var weeks= [7]string{ "MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"}
type Job interface{
	ToString() string
	Run() error
}

type FuncJob struct{
	name string
	f func()
}
func (fj *FuncJob) ToString() string{
	return fj.name
}
func (fj *FuncJob) Run() error{
	fj.f()
	return nil
}


type JobTask struct{
	job Job
}
func (this *JobTask) ToString() string{
	return this.job.ToString()
}
func (this *JobTask) Execute(no int){
	this.job.Run()
}

type Trigger struct{
	Cron string
	cronE *CronParser
	Jobs map[string]Job
}

func NewTrigger(cron string) Trigger{
	tg:= Trigger{cronE:NewCronParser(cron),Jobs:make(map[string]Job)}
	return tg
}

func (this *Trigger) GetNext() time.Time{
	return this.GetNextDate(time.Now())
}
func (this *Trigger) addJob(job Job){
	this.Jobs[job.ToString()]=job
}

func (this *Trigger) GetNextDate(t time.Time) time.Time{
	 checkss,_:=regexp.MatchString("\\d+",this.cronE.ss)
	 checkmi,_:=regexp.MatchString("\\d+",this.cronE.mi)
	 checkhh,_:=regexp.MatchString("\\d+",this.cronE.hh)
	 checkdd,_:=regexp.MatchString("\\d+",this.cronE.dd)
	 checkmm,_:=regexp.MatchString("\\d+",this.cronE.mm)
	 t=t.Add(time.Second)
	 var flag bool
	 for i:=0;i<2<<32;i++{
	 	 t=t.Add(time.Second)
		 flag=this.cronE.Valid(t)
		 if flag{
		 	return t;
		 }
		 if(checkss&&t.Second()==(int)(gosml.ConvertToInt(this.cronE.ss))){
		 	break
		 }
	 }
	 for i:=0;i<2<<32;i++{
	 	 t=t.Add(time.Minute)
		 flag=this.cronE.Valid(t)
		 if flag{
		 	return t;
		 }
		 if(checkmi&&t.Minute()==(int)(go
		 		gosml.ConvertToInt(this.cronE.mi))){
		 	break
		 }
	 }
	 for i:=0;i<2<<32;i++{
	 	 t=t.Add(time.Hour)
		 flag=this.cronE.Valid(t)
		 if flag{
		 	return t;
		 }
		 if(checkhh&&t.Hour()==(int)(gosml.ConvertToInt(this.cronE.hh))){
		 	break
		 }
	 }
	 for i:=0;i<2<<32;i++{
	 	 t=t.AddDate(0,0,1)
		 flag=this.cronE.Valid(t)
		 if flag{
		 	return t;
		 }
		 if(checkdd&&t.Day()==(int)(gosml.ConvertToInt(this.cronE.dd))){
		 	break
		 }
	 }
	 for i:=0;i<2<<32;i++{
	 	 t=t.AddDate(0,1,0)
		 flag=this.cronE.Valid(t)
		 if flag{
		 	return t;
		 }
		 if(checkmm&&(int)(t.Month())==(int)(gosml.ConvertToInt(this.cronE.mm))){
		 	break
		 }
	 }
	 return t
}


type CronParser struct{
	ss,mi,hh,dd,mm,ww,yy string
	ok bool
}

func NewCronParser(cronElp string) *CronParser{
	cp:=&CronParser{"0","0","0","*","*","*","*",false}
	crons:=strings.Split(cronElp," ")
	cp.ok=len(crons)>5
	cp.ss=crons[0];cp.mi=crons[1];cp.hh=crons[2];cp.dd=crons[3];cp.mm=crons[4];cp.ww=crons[5];if len(crons)>6{cp.yy=crons[6]}
	return cp
}

func (cp *CronParser) Valid(t time.Time) bool{
	ss,mi,hh,dd,mm,ww,yy:=t.Second(),t.Minute(),t.Hour(),t.Day(),t.Month(),t.Weekday(),t.Year()
	if(cp.isRight(cp.ss,int64(ss))&&cp.isRight(cp.mi,int64(mi))&&cp.isRight(cp.hh,int64(hh))){
		flag:=cp.isRight(cp.dd,int64(dd))||(strings.EqualFold(cp.dd,"L")&&t.Add(time.Hour*24).Day()==1)
		if flag{
			for i,v:=range months{
				cp.mm=strings.ReplaceAll(cp.mm,v,gosml.ConvertToString(i+1))
			}
			flag=cp.isRight(cp.mm,int64(mm))
			if flag{
				for i,v:=range weeks{
					cp.ww=strings.ReplaceAll(cp.ww,v,gosml.ConvertToString(i+1))
				}
				wd:=getWeek(int64(ww));
				flag=cp.isRight(cp.ww,wd)
				if !flag{
					flag=wd==gosml.ConvertToInt(gosml.SubStr(cp.ww,0,1))
					flag=flag&&strings.Contains(cp.ww,"L")&&dd>t.AddDate(0,0,7).Day()
				}
				return flag&&cp.isRight(cp.yy,int64(yy))
			}
			return flag
		}
	}
	return false
}
func getWeek(ww int64) int64{
	if ww==0{
		return 7
	}else{
		return ww
	}
}
func (cp *CronParser) isRight(evalStr string,ct int64) bool{
		moreR,_:=regexp.MatchString("\\d+/\\d+",evalStr)
		if !moreR {
			if !strings.Contains(evalStr,"-"){ 
				return evalStr=="*"||evalStr=="?"||contains(strings.Split(evalStr,","),strconv.FormatInt(ct,10))
			}else{
				ss:=strings.Split(evalStr,"-")
				return gosml.ConvertToInt(ss[1])<=ct&&sml.ConvertToInt(ss[1])>=ct
			}
		}else{
			ss:=strings.Split(evalStr,"/")
			start,limit:=gosml.ConvertToInt(ss[1]),gosml.ConvertToInt(ss[1])
			return ct-start>=0 &&(ct-start)%limit==0
		}
}
func contains(strs []string,e string) bool{
	for _,v:=range strs{
		if(v==e){	
			return true
		}
	}
	return false
}


//func (t *Trigger) 


type Scheduler struct{
	Triggers map[string]Trigger
	UM sync.RWMutex
	duration  time.Duration
	MQ *ManagedQueue
	stopChan chan bool
}
func (this *Scheduler) SchedulerFuncJob(cron string,name string,f func()){
	this.SchedulerJob(cron,&FuncJob{name:name,f:f})
}
func (this *Scheduler) ExecuteTask(task Task){
	this.MQ.AddTask(task)
}
func (this *Scheduler) SchedulerJob(cron string,job Job) error{
	this.UM.Lock()
	defer this.UM.Unlock()
	var tg Trigger
	flag:=false
	for key,val:=range this.Triggers{
		if key==cron{
			tg=val
			flag=true
		}
	}
	if !flag{
		tg=NewTrigger(cron)
		if !tg.cronE.ok{
			return errors.New("cron error["+cron+"]")
		}
		this.Triggers[cron]=tg
	}else{
		tg=this.Triggers[cron]
	}
	tg.addJob(job)
	return nil
}
func NewDefaultScheduler() *Scheduler{
	return NewScheduler("defaultScheduler",5,10000)
}
func NewScheduler(name string,maxRoutines int,maxDepth int32) *Scheduler{
	scheduler:=&Scheduler{MQ:New(name, maxRoutines, maxDepth),Triggers:make(map[string]Trigger),stopChan:make(chan bool),duration:time.Second}
	scheduler.MQ.Monitor(true)
	return scheduler
}
func (this *Scheduler) WithTicker(d time.Duration) *Scheduler{
	this.duration=d
	return this;
}
func (this *Scheduler) StopScheduler(){
	this.stopChan <- true
}
func (this *Scheduler) InitScheduler(){
	time.Sleep(time.Second-time.Nanosecond*time.Duration(time.Now().UnixNano()%time.Now().Unix()))
	ticker:=time.NewTicker(this.duration)
	go func() {
		defer ticker.Stop()
		for{
			select{
			    case <-ticker.C:{
			        now:=time.Now()
			        for _,trigger:=range this.Triggers{
				        if(trigger.cronE.Valid(now)){
						    for _,job:=range trigger.Jobs{
							     this.MQ.AddTask(&JobTask{job:job})
						    }
				        }
			        }
			    }
			    case stop:=<-this.stopChan:{
				    if(stop){
				    	this.MQ.Destory()
				    	close(this.stopChan)
					    return
				    }
			    }
			}
		}
	}()
}

