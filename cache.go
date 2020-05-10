package gosml

import (
	"sync"
	"time"
	"strings"
)

type CacheManager interface{
	Get(k string) []byte
	Put(k string,v []byte,seconds int64)
	Remove(k string)
	Contain(k string) bool
	Clear() int
	ClearKeyStart(kre string) int
	GetKeyStart(kre string) map[string]interface{}
}

type DefaultCacheManager struct{
	Data map[string][]byte
	Ttls map[string]int64
	Lock sync.RWMutex
}

func NewDefaultCacheManager() *DefaultCacheManager {
   cachemanager:= &DefaultCacheManager{Data:make(map[string][]byte),Ttls:make(map[string]int64)}
   go cachemanager.check()
   return cachemanager
}
func (this *DefaultCacheManager) Get(k string) []byte {
    this.Lock.RLock()
    defer this.Lock.RUnlock()
    if v, exit := this.Data[k]; exit {
        return v
    }
    return nil
}
func (this *DefaultCacheManager) Put(k string,bs []byte,seconds int64){
	this.Lock.Lock()
	defer this.Lock.Unlock()
    this.Data[k] = bs
    if seconds<1{
	    seconds=86400
    }
    this.Ttls[k] =seconds+time.Now().Unix()
}
func (this *DefaultCacheManager) Remove(k string){
	this.Lock.Lock()
	defer this.Lock.Unlock()
	delete(this.Data,k)
	delete(this.Ttls,k)
}
func (this *DefaultCacheManager) Contain(k string) bool{
	this.Lock.RLock()
    defer this.Lock.RUnlock()
	return this.Data[k]!=nil
}
func (this *DefaultCacheManager) Clear() int{
    i:=0
	for k, _ := range this.Data {
		i++
	   this.Remove(k)
	}
	return i
}
func (this *DefaultCacheManager) ClearKeyStart(kre string) int{
    i:=0
	for k, _ := range this.Data {
		if strings.HasPrefix(k, kre){
			i++;
			this.Remove(k)
		}
	}
	return i
}
func (this *DefaultCacheManager) GetKeyStart(kre string) map[string]interface{}{
	result:=make(map[string]interface{})
	for k, v := range this.Data {
		if strings.HasPrefix(k, kre){
			result[k]=v
		}
	}
	return result
}
func (this *DefaultCacheManager) check(){
	for {
		for k,v:=range this.Ttls{
			if v<time.Now().Unix(){
				this.Remove(k)
			}
		}
		time.Sleep(10*time.Second)
	}
}
