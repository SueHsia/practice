package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"rpc"

	"github.com/apache/thrift/lib/go/thrift"
)

func main() {
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	protocolFactory := thrift.NewTCompactProtocolFactory()

	transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "9090"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}

	useTransport, err := transportFactory.GetTransport(transport)
	client := rpc.NewLuckyServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to 127.0.0.1:9090", " ", err)
		os.Exit(1)
	}
	defer transport.Close()

	ctx := context.Background()
	fmt.Println(int64(1))
	res, err := client.DoQuery(ctx, int64(34), "ok", "1234")
	if err != nil {
		log.Println("Echo failed:", err)
		return
	}
	log.Println("response:", res)
	fmt.Println("well done")

	res1, err1 := client.GoodsList(ctx, int64(1), "ok", "111")
	if err1 != nil {
		log.Println("Echo failed:", err1)
		return
	}
	log.Println("response:", res1)
	fmt.Println("well done")
}
