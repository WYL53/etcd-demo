package main

import (
	"context"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func main() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	lease := clientv3.NewLease(cli)

	defer lease.Close()

	leaseGrantResponse, err := lease.Grant(context.TODO(), 10)
	if err != nil {
		panic(err)
	}

	lease.KeepAlive(context.TODO(), leaseGrantResponse.ID)

	key := "/cron/job/lock"

	watcher := clientv3.NewWatcher(cli)
	watcher.Watch(context.TODO(), key, w)

	defer watcher.Close()
}
