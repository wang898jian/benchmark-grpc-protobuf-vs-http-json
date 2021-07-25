package benchmarks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/plutov/benchmark-grpc-protobuf-vs-http-json/http-json"
)

func init() {
	go httpjson.Start()
	time.Sleep(time.Second)
	s, err := strconv.ParseInt(os.Args[len(os.Args)-1], 10, 32)
	if err != nil {
		panic(err)
	}
	scale = int(s)
	fmt.Println("scale:", scale)
	time.Sleep(time.Second)
}

func BenchmarkHTTPJSON(b *testing.B) {
	for i := 0; i <= scale; i++ {
		waitgroup.Add(1)
		client := &http.Client{}
		for n := 0; n < b.N; n++ {
			doPost(client, b)
		}
		waitgroup.Done()
	}
	waitgroup.Wait()

}

func doPost(client *http.Client, b *testing.B) {
	u := &httpjson.User{
		Email:    "foo@bar.com",
		Name:     "Bench",
		Password: "bench",
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(u)

	resp, err := client.Post("http://127.0.0.1:60001/", "application/json", buf)
	if err != nil {
		b.Fatalf("http request failed: %v", err)
	}

	defer resp.Body.Close()

	// We need to parse response to have a fair comparison as gRPC does it
	var target httpjson.Response
	decodeErr := json.NewDecoder(resp.Body).Decode(&target)
	if decodeErr != nil {
		b.Fatalf("unable to decode json: %v", decodeErr)
	}

	if target.Code != 200 || target.User == nil || target.User.ID != "1000000" {
		b.Fatalf("http response is wrong: %v", resp)
	}
}
