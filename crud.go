package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}

	fmt.Println("connect to etcd success")
	defer cli.Close()

	KV := clientv3.NewKV(cli)

	Key := "v1"
	// put
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, err = KV.Put(ctx, Key, "dsb")
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	// get
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	// go func() {
	// 	waitTimeout(cancel)
	// }()
	resp, err := KV.Get(ctx, Key)
	// cancel()

	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}

	//del
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	fmt.Println("start delete")
	delResp, err := KV.Delete(ctx, Key, clientv3.WithPrevKV())
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

func waitTimeout(cancel context.CancelFunc) {
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("太久了，我要 取消操作")
		cancel()
	}
}
