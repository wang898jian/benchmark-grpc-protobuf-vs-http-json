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

func init() {
	go grpcprotobuf.Start()
	s, err := strconv.ParseInt(os.Args[len(os.Args)-1], 10, 32)
	if err != nil {
		panic(err)
	}
	scale = int(s)
	fmt.Println("scale:", scale)
	time.Sleep(time.Second)
}

func BenchmarkGRPCProtobuf(b *testing.B) {
	for i := 0; i <= scale; i++ {
		waitgroup.Add(1)
		go func() {
			conn, err := g.Dial("127.0.0.1:60000", g.WithInsecure())
			if err != nil {
				b.Fatalf("grpc connection failed: %v", err)
			}

			client := proto.NewAPIClient(conn)
			for n := 0; n < b.N; n++ {
				doGRPC(client, b)
			}
		}()
	}
	waitgroup.Wait()
}

func doGRPC(client proto.APIClient, b *testing.B) {
	resp, err := client.CreateUser(context.Background(), &proto.User{
		Email:    "foo@bar.com",
		Name:     "Bench",
		Password: "bench",
	})

	if err != nil {
		b.Fatalf("grpc request failed: %v", err)
	}

	if resp == nil || resp.Code != 200 || resp.User == nil || resp.User.Id != "1000000" {
		b.Fatalf("grpc response is wrong: %v", resp)
	}
}
