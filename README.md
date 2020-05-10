# gosml

 gosml 做为 [![sml](https://github.com/huangwenjimmy/sml) GOLANG实现，保留了大部分工具类（Https,SqlTemplate,ManagedQueue,Cron,Scheduler,gdbc），可基于配置快速开发restful接口
 
 ## Features
 
 ** SqlTemplate 快速方便的数据查询模板，支持配置存放数据库,或本地文件远程文件。
 
 ** Https 链式风格，对各类查询，参数[url参数，表头参数，body参数，文件流]等进行封装方便统一操作
 
 ** Cron cron表达示，实现95%以上，并且增加符合国人的书写如周相关设定
 
 ** ManagedQueue 线程队列管理器，支持指定协程数，队列大小，运行任务监控统计、重新下发等
 
 ** Scheduler 调度器，支持cron表达款任务，秒级精确，轻松支持万级以上任务调度，支持周期调度，单任务调度，延时调度
 
 ** Gdbc 类JDBC操作，支持回调处理，提供常用查询api
 
 ## Usage
 
 ### Gdbc
 
 快速使用 gddc
 
 ```go
 package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/huangwenjimmy/gosml"
	"fmt"
)


func main1(){
	//init sml sqltemplate
	stq:=gosml.Init(gosml.Ds{Id:"defJt",Url:"root:rooxxx@tcp(192.168.1.xxx:3306)/xxxx",DbType:"mysql"})
	stq.Cache(gosml.NewDefaultCacheManager()) //添加缓存
	result_map,_:=stq.Gdbcs["defJt"].QueryForMap("select 1 as a where 1=?",1)
	fmt.Println(result_map) //return map[string]interface{}
	result_list,_:=stq.Gdbcs["defJt"].QueryForList("select * from (select 1 as a union all select 2 as a) t where a=$1",2)
	fmt.Println(result_list)//result []map[string]interface{}
	stq.Gdbcs["defJt"].QueryForCall("select * from sml_if where limit ?",func(row map[string]interface{}) (bool,error){
			//适用于导出等流式操作,返回bool error供上级判断
			return true,nil
		},10000)
	//sml template api支持配置sql参考java-api
	result_multi,_:=stq.Query("test-if",map[string]string{"id":"as","b":"hw","age":"30","time":"2020-04-22T10:06:04Z","dd":"1,2,3","times":"201904,201908"})
	fmt.Println(result_multi)//result    map list page
}
 ```
  ### Cron
  
  crontab表达示实现，支持秒级，下次执行时间，表达示内容
* ss mi HH day-of-mon Mon day-of-w Year[可不选]
*  0 0/5 * * * ?          //从0点，每隔5分钟执行  
*  0 0,30 9-20 * * ?      //每天9点到20点,0分和30分执行
*  0 0 9 1 * ?            //每月1号9点整执行
*  0 40 23 L * ?          //每月最后一天23:40分执行
*  0 0 9 ? * MON          //每周一9点整执行 
*  0 0 9 ? * FRIL         //每月最后一个周期五9点整执行
*  0 0 9 1 JUN-DEC *      //7月至12月每1号9点整执行

```go
package main

import (
	"github.com/huangwenjimmy/gosml"
	"github.com/huangwenjimmy/gosml/queue"
	"fmt"
)

func main(){
	cp:=queue.NewCronParser("0 1 9 * * ?")
	flag:=cp.Valid(gosml.ConvertStrToDate("20200428090100"))
	fmt.Println(flag) //true
	//trigger
	tri:=queue.NewTrigger("0 1 9 * * ?")
	fmt.Println(tri.GetNext())//当前时间下次执行时间
	fmt.Println(tri.GetNextDate(gosml.ConvertStrToDate("20200428090100")))//指定时间下次执行时间
}
```

### ManagedQueue|Scheduler

调度器+线程管理器结合使用，高效流转各内部环节任务管理

```go
package main

import (
	"github.com/huangwenjimmy/gosml/queue"
	"fmt"
	"time"
)

func main(){
	ds:=queue.NewScheduler("defaultScheduler",10,10000)//创建调度器指定10个协程，10000个队列深度
	ds.InitScheduler()//初始化
	//下发周期任务
	ds.SchedulerJob("3/5 * * * * ?",&DJ{"job1"})//添加结构体任务  implements  queue.Job method ToString|Run
	//下发方法任务
	ds.SchedulerFuncJob("0 0/5 * * * ?","funcJob1",func(){
		fmt.Println("funcJob1",time.Now())
	})
	//执行一次性任务
	//ds.ExecuteTask(task)//task implements queue.Task
	time.Sleep(time.Hour)
}
type DJ struct{
	name string
}
func (dj *DJ) ToString() string{
	return dj.name
}
func (dj *DJ) Run() error{
	fmt.Println(dj.name,time.Now())
	return nil
}
```

### Https   
 支持常见操作

```go
package main

import (
	"github.com/huangwenjimmy/gosml"
	"fmt"
	"time"
	"strings"
)

func main(){
	//添加表头  basic 认证   指定连接读写超时时间
	https:=gosml.NewGetHttps("https://localhost:16003/bus/sml/query/if-cfg-interMngLike").
	Header("content-type","application/json").
	BasicAuth("user:passwd").
	ConnectTimeout(time.Second).RWTimeout(time.Second).
	Param("ids","10001").Param("name","测试")
	https.Execute()
	//支持对返回结果进行多种操作,获取流，反序列化，下载等
	fmt.Println(https.GetBodyString())//https.GetBodyBytes()|https.GetBodyTo(writer)|https.GetBodyToFile(*File)
	//post-body请求，body支持   string []byte  io.reader
	https=gosml.NewPostBodyHttps("https://localhost:16003/bus/sml/query/if-cfg-interMngLike").Body(`{"ids":"1001","name":"测试"}`)
	https.Execute()
	fmt.Println(https.GetBodyString())
	//post-upload上传  需要指定Multipart  表单参数Param url参数请自行拼接至url  通过gosml.UpFile  可添加多个文件
	https=gosml.NewPostFormHttps("https://localhost:16003/bus/sml/web/file/upload").Multipart().
	Param("filepath","../logs").
	UpFile(&gosml.UpFile{FormName:"file",FileName:"hello1.txt",Input:strings.NewReader("helloworld234天啊\n我是上传的内容")})
	https.Execute()
	fmt.Println(https.Response.Status,https.GetBodyString())
}
```



 
  
 