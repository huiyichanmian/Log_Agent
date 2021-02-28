package taillog

import (
	"Log_Agent/etcd"
	"Log_Agent/kafka"
	"context"
	"fmt"
	"github.com/hpcloud/tail"
)

// 专门从日志文件收集日志模块

// TailTask：一个日志收集的任务
type TailTask struct{
	path string
	topic string
	instance *tail.Tail
	// 为了能够实现退出t.run()
	ctx context.Context
	cancelFunc context.CancelFunc
}


func NewTailTask(path,topic string)(tailObj *TailTask){
	ctx,cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path: path,
		topic:topic,
		ctx:ctx,
		cancelFunc: cancel,
	}
	tailObj.init()  // 根据路径去打开对应的日志
	return
}

func (t *TailTask)init(){
	config := tail.Config{
		ReOpen: true,  // 重新打开
		Follow: true,  // 是否跟随
		Location: &tail.SeekInfo{Offset: 0,Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,   // 文件不存在不报错
		Poll: true,   // 轮询
	}
	var err error
	if t.instance, err = tail.TailFile(t.path,config); err != nil {
		fmt.Printf("Init taillog failed,err:%v\n",err)
	}
	// 当goroutine执行的函数推出的时候，goroutine就结束了
	go t.Run()  // 直接去采集日志发送到kafka
}

func(t *TailTask)ReadChan() <-chan *tail.Line {
	return t.instance.Lines
}

// 给kafka发送消息
func (t *TailTask) Run(){
	for {
		select {
		case <-t.ctx.Done():
			fmt.Printf("tail task:%s_%s 结束了...\n",t.path,t.topic)
			return
		case line := <- t.ReadChan():
			// 3.2给kafka发送消息(同步等待)
			// kafka.SendTokafka(t.topic,line.Text)
			// note: 优化
			//先把日志数据发送到一个通道, kafka那个包中有单独的goroutine去取日志数据发到kafka
			kafka.SendToChan(t.topic,line.Text)
		}
	}
}

// 向外暴露一个函数，向taskMgr的newConfChan
func NewConfChan() chan <- []*etcd.LogEntry{
	return taskMgr.newConfChan
}


//func Init(fileName string)(err error){
//	config := tail.Config{
//		ReOpen: true,  // 当文件大于设定的值后重新打开
//		Follow: true,  // 是否跟随
//		Location:&tail.SeekInfo{Offset: 0, Whence: 2},  // 从文件哪个地方读
//		MustExist: false,  // 允许文件不存在
//		Poll: true,  // 轮询
//	}
//
//	// 打开文件，开始读取数据
//	tailObj, err = tail.TailFile(fileName,config)
//	if err != nil{
//		fmt.Printf("Init taillog failed,err:%v\n",err)
//		return err
//	}
//	return
//}
//
//
//func ReadChan() <-chan *tail.Line{
//	return tailObj.Lines
//}
