package taillog

import (
	"Log_Agent/etcd"
	"fmt"
	"time"
)

var taskMgr *tailLogMgr

// tailTask 管理者
type tailLogMgr struct{
	logEntry []*etcd.LogEntry
	taskMap map[string]*TailTask
	newConfChan chan []*etcd.LogEntry
}

func Init(logEntryConf []*etcd.LogEntry){
	taskMgr = &tailLogMgr{
		logEntry: logEntryConf,
		taskMap: make(map[string]*TailTask,16),
		newConfChan: make(chan []*etcd.LogEntry),  // 无缓冲的通道，没有值一直阻塞
	}
	for _, logEntry := range logEntryConf{
		// etcd.LogEntry
		// logEntry.Path  要收集日志的路劲
		// 初始化tail
		// 初始化的时候起了多少个tailtask,都要记下来，为了后续判断的方便
		tailObj := NewTailTask(logEntry.Path,logEntry.Topic)
		mk := fmt.Sprintf("%s_%s",logEntry.Path,logEntry.Topic)  // 拼接map的key
		taskMgr.taskMap[mk] = tailObj
	}
	go taskMgr.run()
}


// 监听自己的newConfChan,有了新的配置过来之后就做对应的处理
func (t *tailLogMgr) run (){
	for{
		select {
		case newConf := <- t.newConfChan:
			for _, conf := range newConf{
				mk := fmt.Sprintf("%s_%s",conf.Path,conf.Topic)
				if _, ok := t.taskMap[mk];ok{
					// 原来就有,不需要操作
					fmt.Println("发现原来就有，不需要在增加了")
					continue
				}
				// 1.配置新增
				fmt.Println("新增配置")
				tailObj := NewTailTask(conf.Path,conf.Topic)
				t.taskMap[mk] = tailObj
			}
			// 找出原来t.logEntry有，但是newConf中没有的，需要删除
			for _,c1 := range t.logEntry{
				isDelete := true
				for _,c2 := range newConf{
					if c2.Path == c1.Path && c2.Topic == c1.Topic{
						isDelete = false
						continue
					}
				}
				if isDelete{
					// 把c1对应的tailObj给停掉
					mk := fmt.Sprintf("%s_%s",c1.Path,c1.Topic)
					// t.taskMap[mk] ==> tailObj
					t.taskMap[mk].cancelFunc()
				}
			}
			// 2.配置删除
			// 3.配置变更
			fmt.Println("新的配置来了",newConf)
		default:
			time.Sleep(time.Second)
		}
	}
}