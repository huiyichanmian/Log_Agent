package main

import (
	"Log_Agent/conf"
	"Log_Agent/etcd"
	"Log_Agent/kafka"
	"Log_Agent/taillog"
	"Log_Agent/utils"
	"fmt"
	"gopkg.in/ini.v1"
	"sync"
	"time"
)

var (
	cfg = new(conf.AppConf)
)


// logAgent入口程序
func main() {
	// 0. 加载配置文件
	if err := ini.MapTo(cfg, "./conf/config.ini"); err != nil {
		fmt.Printf("load ini failed, err:%v\n", err)
		return
	}
	// 1. 初始化kafka连接
	if err := kafka.Init([]string{cfg.KafkaConf.Address},cfg.ChanMaxSize); err != nil {
		fmt.Printf("Init kafka failed, err:%v\n", err)
		return
	}

	// 2.初始化ETCD
	if err := etcd.Init(cfg.EtcdConf.Address,time.Duration(cfg.EtcdConf.Timeout)*time.Second);err!=nil{
		fmt.Printf("Init etcd failed,err:%v\n",err)
		return
	}
	// 为了实现每个logAgent都拉取自己独有的配置，所以要以自己的ip地址来区分
	ipStr, err := utils.GetOutboundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(cfg.EtcdConf.Key,ipStr)
	// 2.1 从etcd中获取日志收集项的信息
	logEntryConf, err := etcd.GetConf(etcdConfKey)
	if err != nil{
		fmt.Printf("etcd.GetConf failed,err:%v\n",err)
		return
	}
	fmt.Printf("get conf from etcd success,%v\n",logEntryConf)

	for index, value := range logEntryConf{
		fmt.Printf("index:%v,value:%v\n",index,value)
	}

	// 3.收集日志，发给kafka
	// 3.1 循环每一个日志收集项，创建TailObj
	taillog.Init(logEntryConf)
	// 3.2 派一个哨兵去见识日志收集项的变化（有变化及时通知我的logAgent实现热加载配置）
	// 因为NewConfChan访问了taskMgr的newConfChan,这个channel是在taillog.Init(logEntryConf)执行的初始化
	newConfChan := taillog.NewConfChan()  // 从taillog包中获取对外暴露的通道
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(etcdConfKey,newConfChan)  // 哨兵发现最新的配置信息会通知上面的那个通道
	wg.Wait()

	//// 2. 打开日志文件准备收集日志
	//if err := taillog.Init(cfg.FileName); err != nil {
	//	fmt.Printf("Init taillog failed, err:%v\n", err)
	//	return
	//}
	//run()
}
