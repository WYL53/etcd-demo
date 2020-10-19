package main

import(
	"github.com/coreos/go-etcd/etcd"
)

func main()  {
	machines := []string{"http://127.0.0.1:2379"}
    client := etcd.NewClient(machines)

    if _, err := client.Set("/foo", "bar", 0); err != nil {
        log.Fatal(err)
    }
}