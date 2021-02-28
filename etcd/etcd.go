package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	cli *clientv3.Client
)

type LogEntry struct {
	Path string `json:"path"`  // 日志存放的路径
	Topic string `json:"topic"`  // 日志要发往Kafka中的哪个Topic
}


// 初始化etcd的函数
func Init(addr string, timeout time.Duration)(err error){
	if cli,err = clientv3.New(clientv3.Config{
		Endpoints: []string{addr},
		DialTimeout: timeout,
	}); err != nil{
		fmt.Printf("connect to etcd failed,err:%v\n",err)
		return
	}
	return
}

// 从ETCD中根据KEY获取配置项
func GetConf(key string)(logEntryConf []*LogEntry, err error){
	// get
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	resp, err := cli.Get(ctx,key)
	cancel()
	if err != nil{
		fmt.Printf("get from etcd failed, err:%v\n",err)
		return
	}
	for _, ev := range resp.Kvs{
		//fmt.Printf("%s:%s\n",ev.Key,ev.Value)
		if err = json.Unmarshal(ev.Value, &logEntryConf); err != nil{
			fmt.Printf("unmarshal etcd value failed, err:%v\n",err)
			return
		}
	}
	return
}


func WatchConf(key string, newConfCh chan<-[]*LogEntry){
	ch := cli.Watch(context.Background(),key)
	for wresp := range ch{
		for _, evt := range wresp.Events{
			fmt.Printf("Type:%v key:%v value:%v\n",evt.Type,string(evt.Kv.Key),string(evt.Kv.Value))
			// 通知taillog.taskMgr
			// 1.先判断操作的类型
			var newConf []*LogEntry
			// 如果是删除操作，手动传递一个空的配置项
			if evt.Type != clientv3.EventTypeDelete{
				if err := json.Unmarshal(evt.Kv.Value,&newConf); err != nil {
					fmt.Printf("unmarshal failed,err:%v\n",err)
					continue
				}
			}
			fmt.Printf("get new fonf:%v\n", newConf)
			newConfCh <- newConf
		}
	}
}

