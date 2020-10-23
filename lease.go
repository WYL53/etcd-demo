package main

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var (
	cli                        *clientv3.Client
	err                        error
	kv                         clientv3.KV
	lease                      clientv3.Lease
	leaseGrantResponse         *clientv3.LeaseGrantResponse
	putResponse                *clientv3.PutResponse
	getResponse                *clientv3.GetResponse
	leaseKeepAliveResponse     *clientv3.LeaseKeepAliveResponse
	leaseKeepAliveResponseChan <-chan *clientv3.LeaseKeepAliveResponse
)

func main() {

	config := clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}

	if cli, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	defer cli.Close()

	kv = clientv3.NewKV(cli)

	lease := clientv3.NewLease(cli)

	if leaseGrantResponse, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	}

	leaseId := leaseGrantResponse.ID
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if leaseKeepAliveResponseChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println("KeepAlive err : ", err)
		return
	}

	if putResponse, err = kv.Put(context.TODO(), "/cron/job1/lock", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println("put key error : ", err)
		return
	}
	fmt.Println("put success ", putResponse.Header.Revision)

	go func() {

		for {
			select {
			case leaseKeepAliveResponse = <-leaseKeepAliveResponseChan:

				if leaseGrantResponse == nil {
					fmt.Println("续约失败")
					break
				}

				fmt.Println("续约成功 ：", leaseKeepAliveResponse.ID, " ttl : ", leaseKeepAliveResponse.TTL)
			}
		}

	}()

	for {
		if getResponse, err = kv.Get(context.TODO(), "/cron/job1/lock"); err != nil {
			fmt.Println("获取 key失败：", err)
			return
		}
		if getResponse.Count == 0 {
			fmt.Println("租约已过期了")
			break
		}

		fmt.Println("租约还没到期 :", getResponse.Kvs[0])
		time.Sleep(time.Second * 2)
	}

}
