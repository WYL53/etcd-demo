package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"118.89.22.227:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}

	fmt.Println("connect to etcd success")
	defer cli.Close()

	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	Key := "q1mi"
	_, err = cli.Put(ctx, Key, "dsb")
	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, Key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}

	//del
	fmt.Println("start delete")
	delResp, err := cli.Delete(context.TODO(), Key, clientv3.WithPrevKV())
	if err != nil {
		fmt.Printf("del key  failed , err： %v \n", err)
		return
	}
	if delResp.Deleted > 0 {
		for _, kv := range delResp.PrevKvs {
			fmt.Println(string(kv.Key), string(kv.Value))
		}
	}

}
