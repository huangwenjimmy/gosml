package queue

import (
	"sort"
	"sync"
	"time"
	"fmt"
)

type  Delayer interface{
	GetId() string
	GetDelay() int64
}

type DelayQueue struct{
	delayers []Delayer
	UM sync.RWMutex
	dequeueChan chan Delayer
	stopChan chan bool
}

func (dq *DelayQueue) Len() int{
	return len(dq.delayers)
}
func (dq *DelayQueue) Less(i, j int) bool {
	return dq.delayers[i].GetDelay() < dq.delayers[j].GetDelay()
}
func (dq *DelayQueue) Swap(i, j int) {
	dq.delayers[i],dq.delayers[j] = dq.delayers[j],dq.delayers[i]
}

func NewDelayQueue(maxDepth int) *DelayQueue{
	dq:= &DelayQueue{delayers:make([]Delayer,0),dequeueChan:make(chan Delayer,maxDepth),stopChan:make(chan bool)}
	dq.initDelayQueue()
	return dq
}
func (dq *DelayQueue) AddDelayer(delayer Delayer){
	dq.UM.Lock()
	defer dq.UM.Unlock()
	dq.delayers=append(dq.delayers,delayer)
	sort.Sort(dq)
}
func (dq *DelayQueue) sliceDelayer(j int){
	dq.UM.Lock()
	defer dq.UM.Unlock()
	dq.delayers=dq.delayers[j:]
}
func (dq *DelayQueue) DeQueue() Delayer{
	delayer:=<-dq.dequeueChan
	return delayer
}
func (dq *DelayQueue) Remove(delayerid string) Delayer{
	dq.UM.Unlock()
	defer dq.UM.Unlock()
	for i,delayer:=range dq.delayers{
		if(delayerid==delayer.GetId()){
			dq.delayers=append(dq.delayers[:i],dq.delayers[i+1:]...)
			return delayer
		}
	}
	return nil
}
func (dq *DelayQueue) initDelayQueue(){
	ticker:=time.NewTicker(time.Second)
	go func(){
		defer ticker.Stop()
		for{
			select{
				case <-ticker.C:{
					j:=0
					for _,v:=range dq.delayers{
						if(v.GetDelay()<=0){
							dq.dequeueChan <- v
							j++
						}
					}
					if j>0{
						dq.sliceDelayer(j)
					}
				}
				case stop:=<-dq.stopChan:{
					if(stop){
						return
					}
				}
			}
		}
	}()
}

func debug(args ... interface{}){
	fmt.Println(args...)
}

