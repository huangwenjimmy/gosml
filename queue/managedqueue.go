package queue

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var ErrCapecity=errors.New("capacity is full")

type Task interface{
	Execute(routineNo int)
	ToString() string
}

type delegateTask struct{
	task Task
	resultChannel chan error
}
type executingTasks struct{
	taskMap map[string]time.Time
	um sync.RWMutex
}
func (et *executingTasks) put(taskName string){
	et.um.Lock()
	defer et.um.Unlock()
	et.taskMap[taskName]=time.Now()
}
func (et *executingTasks) remove(taskName string){
	et.um.Lock()
	defer et.um.Unlock()
	delete(et.taskMap,taskName)
}
func (et *executingTasks) contain(taskName string) bool{
	et.um.RLock()
	defer et.um.RUnlock()
	for k,_:=range et.taskMap{
		if(k==taskName){
			return true
		}
	}
	return false
}

type ManagedQueue struct{
	name string
	isSingleTask bool
	executingTasks *executingTasks
	shutdownMasterChannel chan string     //  master routine.
	shutdownSlaveChannel  chan struct{}   //  slave routines.
	shutdownWaitGroup    sync.WaitGroup  // The WaitGroup for shutting down existing routines.
	delegateTaskChannel         chan delegateTask   // Channel used to sync access to the queue.
	taskChannel          chan Task       
	queuedSize           int32          
	activeRoutines       int32           
	maxDepth        int32           
}

func New(name string,maxRoutines int, maxDepth int32) *ManagedQueue{
	managedQueue:=ManagedQueue{
		name :name,
		shutdownMasterChannel: make(chan string),
		shutdownSlaveChannel:  make(chan struct{}),
		delegateTaskChannel:  make(chan delegateTask),
		taskChannel:          make(chan Task, maxDepth),
		queuedSize:           0,
		activeRoutines:       0,
		maxDepth:        maxDepth,
	}
	managedQueue.shutdownWaitGroup.Add(maxRoutines)
	for routineNo := 0; routineNo < maxRoutines; routineNo++ {
		go managedQueue.taskRoutine(routineNo)
	}
	go managedQueue.queueRoutine()
	
	return &managedQueue
}
func  (managedQueue *ManagedQueue) Monitor(isSingleTask bool) *ManagedQueue{
	managedQueue.executingTasks=&executingTasks{taskMap:make(map[string]time.Time)}
	managedQueue.isSingleTask=isSingleTask
	return managedQueue
}
func (managedQueue *ManagedQueue) GetExecutingTasks() map[string]time.Time{
	return managedQueue.executingTasks.taskMap
}


func (managedQueue *ManagedQueue) taskRoutine(routineNo int) {
	for {
		select {
			case <-managedQueue.shutdownSlaveChannel:
				writeStdout(fmt.Sprintf("routineNo %d", routineNo), "slaveRoutine", "Going Down")
				managedQueue.shutdownWaitGroup.Done()
				return
	
			case task := <-managedQueue.taskChannel:
				managedQueue.run(routineNo, task)
				break
		}
	}
}
func (managedQueue *ManagedQueue) run(routineNo int, task Task) {
	defer catchPanic(nil, "routineNo", "run")
	defer atomic.AddInt32(&managedQueue.activeRoutines, -1)
	if	managedQueue.executingTasks!=nil{
		defer managedQueue.executingTasks.remove(task.ToString())
	}
	atomic.AddInt32(&managedQueue.queuedSize, -1)
	atomic.AddInt32(&managedQueue.activeRoutines, 1)
	task.Execute(routineNo)
}
func (managedQueue *ManagedQueue) queueRoutine() {
	for {
		select {
		case <-managedQueue.shutdownMasterChannel:
			writeStdout("Queue", "MasterRoutine", "Going Down")
			managedQueue.shutdownMasterChannel <- "Down"
			return
		case delegateTask := <-managedQueue.delegateTaskChannel:
			if managedQueue.QueuedSize() == managedQueue.maxDepth {
				delegateTask.resultChannel <- ErrCapecity
				continue
			}
			atomic.AddInt32(&managedQueue.queuedSize, 1)
			if	managedQueue.executingTasks!=nil{
				managedQueue.executingTasks.put(delegateTask.task.ToString())
			}
			managedQueue.taskChannel <- delegateTask.task
			delegateTask.resultChannel <- nil
			break
		}
	}
}
func (managedQueue *ManagedQueue) AddTask(task Task) (err error) {
	if managedQueue.isSingleTask&&managedQueue.executingTasks.contain(task.ToString()){
		return nil
	}else{
		return managedQueue.addTask0(task)
	}
}

func (managedQueue *ManagedQueue) addTask0(task Task) (err error){
	defer catchPanic(&err, managedQueue.name, "Add")
	delegateTask := delegateTask{task, make(chan error)}
	defer close(delegateTask.resultChannel)
	managedQueue.delegateTaskChannel <- delegateTask
	err = <-delegateTask.resultChannel
	return err
}

func (managedQueue *ManagedQueue) Destory() (err error) {
	defer catchPanic(&err, managedQueue.name, "Destory")
	writeStdout(managedQueue.name, "Destory", "Started")
	writeStdout(managedQueue.name, "Destory", "Master Routine")
	managedQueue.shutdownMasterChannel <- "Down"
	<-managedQueue.shutdownMasterChannel
	close(managedQueue.delegateTaskChannel)//不接收任务
	close(managedQueue.shutdownMasterChannel)//分发队列停止任务
	writeStdout(managedQueue.name, "Destory", "SLAVER Routines")
	close(managedQueue.shutdownSlaveChannel)//关闭线程不接收任务
	managedQueue.shutdownWaitGroup.Wait()//等待所有任务执行完
	close(managedQueue.taskChannel)//关闭任务通道
	writeStdout(managedQueue.name, "Destory", "Completed")
	return err
}
func (managedQueue *ManagedQueue) GetName() string {
	return managedQueue.name
}

func (managedQueue *ManagedQueue) QueuedSize() int32 {
	return atomic.AddInt32(&managedQueue.queuedSize, 0)
}

func (managedQueue *ManagedQueue) ActiveRoutines() int32 {
	return atomic.AddInt32(&managedQueue.activeRoutines, 0)
}

func catchPanic(err *error, goRoutine string, functionName string) {
	if r := recover(); r != nil {
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)
		writeStdoutf(goRoutine, functionName, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))
		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

func writeStdout(goRoutine string, functionName string, message string) {
	log.Printf("%s : %s : %s\n", goRoutine, functionName, message)
}

func writeStdoutf(goRoutine string, functionName string, format string, a ...interface{}) {
	writeStdout(goRoutine, functionName, fmt.Sprintf(format, a...))
}
