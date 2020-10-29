package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {

	conf := clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(conf)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	key := "/corn/job1/lock"

	go func() {
		kv := clientv3.NewKV(client)
		// ctx , _ := context.WithTimeout(context.Background(),time.Minute)
		after := time.After(time.Minute)
		for {
			select {
			case <-after:
				return
			default:
				kv.Put(context.TODO(), key, "")
				time.Sleep(time.Second)
				kv.Delete(context.TODO(), key)
			}
		}
	}()

	watcher := clientv3.NewWatcher(client)
	defer watcher.Close()

	watchChan := watcher.Watch(context.TODO(), key)
	for {
		select {
		case wc := <-watchChan:
			fmt.Println("收到watch消息：", wc.Events[0].Type, wc.Header.Revision)
		case <-time.After(time.Second * 5):
			goto loop
		}
	}
loop:
}
