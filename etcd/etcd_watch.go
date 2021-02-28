package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed,err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()

	// watch
	ch := cli.Watch(context.Background(), "Negan") // 派一个哨兵一直监视Egan这个key的变化（新增修改删除）
	// 从通道中尝试取值（监视的信息）
	for wresp := range ch {
		for _, evt := range wresp.Events {
			fmt.Printf("Type:%v. key:%v, value:%v\n", evt.Type, string(evt.Kv.Key), string(evt.Kv.Value))
		}
	}
}
