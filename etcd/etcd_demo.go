package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func etcdDemo() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: 5*time.Second,
	})
	if err != nil{
		fmt.Printf("connect to etcd failed,err:%v\n",err)
		return
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()

	// put
	ctx, cancel := context.WithTimeout(context.Background(),time.Second)
	value := `[{"path":"d:/xxx/mysql.log","topic":"mysql_log"}]`
	_, err = cli.Put(ctx,"/logagent/192.168.2.104/collect_config",value)
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed,err:%v\n",err)
		return
	}

	// get
	//ctx, cancel = context.WithTimeout(context.Background(),time.Second)
	//resp,err:=cli.Get(ctx,"Negan")
	//cancel()
	//if err != nil{
	//	fmt.Printf("get from etcd failed,err:%v\n",err)
	//	return
	//}
	//for _, v := range resp.Kvs{
	//	fmt.Printf("%s:%s\n",v.Key,v.Value)
	//}

}