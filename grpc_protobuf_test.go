package benchmarks

import (
	"fmt"
	g "google.golang.org/grpc"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/plutov/benchmark-grpc-protobuf-vs-http-json/grpc-protobuf"
	"github.com/plutov/benchmark-grpc-protobuf-vs-http-json/grpc-protobuf/proto"
	"golang.org/x/net/context"
)

var scale int
var waitgroup sync.WaitGroup
var packageCount int

func init() {
	go grpcprotobuf.Start()
	s, err := strconv.ParseInt(os.Args[len(os.Args)-1], 10, 32)
	if err != nil {
		panic(err)
	}
	pc, err := strconv.ParseInt(os.Args[len(os.Args)-2], 10, 32)
	if err != nil {
		panic(err)
	}
	packageCount = int(pc)
	scale = int(s)
	fmt.Println("grpc scale:", scale)
	fmt.Println("grpc packageCount:", packageCount)
	time.Sleep(time.Second)
}

func TestBenchmarkGRPCProtobuf(t *testing.T) {
	timeBegin := time.Now()
	var timeRun time.Duration
	for i := 0; i <= scale; i++ {
		waitgroup.Add(1)
		go func(timeRun *time.Duration) {
			tmpTime := time.Now()
			conn, err := g.Dial("127.0.0.1:60000", g.WithInsecure())
			if err != nil {
				t.Fatalf("grpc connection failed: %v", err)
			}

			client := proto.NewAPIClient(conn)
			for n := 0; n < packageCount; n++ {
				doGRPC(client, t)
			}
			waitgroup.Done()
			Duration := time.Since(tmpTime)
			fmt.Printf("grpc Parse time is %+v\n", Duration)
		}(&timeRun)
	}
	Duration := time.Since(timeBegin)
	fmt.Printf("grpc time is %+v\n", Duration)
	waitgroup.Wait()
}

func doGRPC(client proto.APIClient, t *testing.T) {
	resp, err := client.CreateUser(context.Background(), &proto.User{
		Email:    "foo@bar.com",
		Name:     "Bench",
		Password: "bench",
	})

	if err != nil {
		t.Fatalf("grpc request failed: %v", err)
	}

	if resp == nil || resp.Code != 200 || resp.User == nil || resp.User.Id != "1000000" {
		t.Fatalf("grpc response is wrong: %v", resp)
	}
}
