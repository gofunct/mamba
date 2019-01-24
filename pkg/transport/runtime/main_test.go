package runtime_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofunct/mamba/pkg/transport/runtime"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofunct/mamba/pkg/transport/testing/api"
	"github.com/gofunct/mamba/pkg/transport/testing/app/server"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
)

var (
	waitForServer = func() { time.Sleep(1 * time.Second) }
)

func orDie(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func startServer(t *testing.T, s *runtime.Engine) func() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.Serve(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			t.Errorf("Engine.Serve returned an error: %v", err)
		}
	}()
	waitForServer()
	return func() {
		s.Shutdown()
		wg.Wait()
	}
}

func Test_server_onlyGateway(t *testing.T) {
	var port int64 = 15261
	s := runtime.New(
		runtime.WithGatewayAddr("tcp", ":"+strconv.FormatInt(port, 10)),
		runtime.WithServers(
			server.NewLibraryServiceServer(),
		),
	)

	defer startServer(t, s)()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/books", port))
	orDie(t, err)
	defer resp.Body.Close()

	if got, want := resp.StatusCode, 200; got != want {
		t.Errorf("Response status is %d, want %d", got, want)
	}

	data, err := ioutil.ReadAll(resp.Body)
	orDie(t, err)

	got := map[string]interface{}{}
	orDie(t, json.Unmarshal(data, &got))
	want := map[string]interface{}{
		"books": []interface{}{
			map[string]interface{}{"book_id": "The Go Programming Language"},
			map[string]interface{}{"book_id": "Programming Ruby"},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Received body differs: (-got +want)\n%s", diff)
	}
}

func Test_server_samePort(t *testing.T) {
	var port int64 = 15261
	addr := ":" + strconv.FormatInt(port, 10)
	s := runtime.New(
		runtime.WithAddr("tcp", addr),
		runtime.WithServers(
			server.NewLibraryServiceServer(),
		),
	)

	defer startServer(t, s)()

	t.Run("http", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/books", port))
		orDie(t, err)
		defer resp.Body.Close()

		if got, want := resp.StatusCode, 200; got != want {
			t.Errorf("Response status is %d, want %d", got, want)
		}

		data, err := ioutil.ReadAll(resp.Body)
		orDie(t, err)

		got := map[string]interface{}{}
		orDie(t, json.Unmarshal(data, &got))
		want := map[string]interface{}{
			"books": []interface{}{
				map[string]interface{}{"book_id": "The Go Programming Language"},
				map[string]interface{}{"book_id": "Programming Ruby"},
			},
		}

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Received body differs: (-got +want)\n%s", diff)
		}
	})

	t.Run("gRPC", func(t *testing.T) {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		orDie(t, err)
		defer conn.Close()

		cli := api_pb.NewLibraryServiceClient(conn)
		resp, err := cli.ListBooks(context.Background(), &api_pb.ListBooksRequest{})
		orDie(t, err)

		want := &api_pb.ListBooksResponse{
			Books: []*api_pb.Book{
				{BookId: "The Go Programming Language"},
				{BookId: "Programming Ruby"},
			},
		}
		if diff := cmp.Diff(resp, want); diff != "" {
			t.Errorf("Received body differs: (-got +want)\n%s", diff)
		}
	})
}

func Test_server_differentPort(t *testing.T) {
	var (
		grpcPort int64 = 15261
		httpPort int64 = 15262
	)

	grpcAddr := ":" + strconv.FormatInt(grpcPort, 10)
	httpAddr := ":" + strconv.FormatInt(httpPort, 10)

	s := runtime.New(
		runtime.WithGrpcAddr("tcp", grpcAddr),
		runtime.WithGatewayAddr("tcp", httpAddr),
		runtime.WithServers(
			server.NewLibraryServiceServer(),
		),
	)

	defer startServer(t, s)()

	t.Run("http", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/books", httpPort))
		orDie(t, err)
		defer resp.Body.Close()

		if got, want := resp.StatusCode, 200; got != want {
			t.Errorf("Response status is %d, want %d", got, want)
		}

		data, err := ioutil.ReadAll(resp.Body)
		orDie(t, err)

		got := map[string]interface{}{}
		orDie(t, json.Unmarshal(data, &got))
		want := map[string]interface{}{
			"books": []interface{}{
				map[string]interface{}{"book_id": "The Go Programming Language"},
				map[string]interface{}{"book_id": "Programming Ruby"},
			},
		}

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Received body differs: (-got +want)\n%s", diff)
		}
	})

	t.Run("gRPC", func(t *testing.T) {
		conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
		orDie(t, err)
		defer conn.Close()

		cli := api_pb.NewLibraryServiceClient(conn)
		resp, err := cli.ListBooks(context.Background(), &api_pb.ListBooksRequest{})
		orDie(t, err)

		want := &api_pb.ListBooksResponse{
			Books: []*api_pb.Book{
				{BookId: "The Go Programming Language"},
				{BookId: "Programming Ruby"},
			},
		}
		if diff := cmp.Diff(resp, want); diff != "" {
			t.Errorf("Received body differs: (-got +want)\n%s", diff)
		}
	})
}
