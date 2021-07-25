package benchmarks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/plutov/benchmark-grpc-protobuf-vs-http-json/http-json"
)

var sc int
var wg sync.WaitGroup
var pc int

func init() {
	go httpjson.Start()
	time.Sleep(time.Second)
	s, err := strconv.ParseInt(os.Args[len(os.Args)-1], 10, 32)
	if err != nil {
		panic(err)
	}
	sc = int(s)
	packageCount, err := strconv.ParseInt(os.Args[len(os.Args)-2], 10, 32)
	if err != nil {
		panic(err)
	}
	pc = int(packageCount)
	fmt.Println("httpjson scale:", sc)
	fmt.Println("httpjson packageCount:", packageCount)
	time.Sleep(time.Second)
}

//func BenchmarkHTTPJSON(b *testing.B) {
//	for i := 0; i <= scale; i++ {
//		waitgroup.Add(1)
//		client := &http.Client{}
//		for n := 0; n < packageCount; n++ {
//			doPost(client, b)
//		}
//		waitgroup.Done()
//	}
//	waitgroup.Wait()
//
//}

func TestBenchmarkHTTPJSON(t *testing.T) {
	timeBegin := time.Now()
	var timeRun time.Duration
	client := &http.Client{}
	for i := 0; i <= sc; i++ {
		wg.Add(1)
		go func(timeRun *time.Duration) {
			tmpTime := time.Now()
			for n := 0; n < pc; n++ {
				doPost(client, t)
			}
			wg.Done()
			Duration := time.Since(tmpTime)
			fmt.Printf("HTTPJSON total parse time is %+v, avg time is %+vns\n", Duration, int(Duration)/pc)
		}(&timeRun)
	}
	Duration := time.Since(timeBegin)
	fmt.Printf("HTTPJSON main time is %+v\n", Duration)
	wg.Wait()
}

func doPost(client *http.Client, b *testing.T) {
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
