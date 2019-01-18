package function

import (
	"bufio"
	"context"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"gocloud.dev/requestlog"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type HandlerFunc func(ctx context.Context, request interface{}) (response interface{}, err error)
type Middleware func(HandlerFunc) HandlerFunc
type HttpHandlerFunc func(http.ResponseWriter, *http.Request)
type RoundTripperFunc func(*http.Request) (*http.Response, error)
type ServeFunc func(addr string, h http.Handler) error
type ShutdownFunc func(ctx context.Context) error
type HijackerFunc func() (net.Conn, *bufio.ReadWriter, error)
type RouterWalkFunc func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error
type RouterBuildFunc func(map[string]string) map[string]string
type UnaryInterceptorFunc func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
type StreamInterceptorFunc func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
type BaseFunc func()
type StringerFunc func() string
type OnErrFunc func(error)
type ErrFunc func() string
type AuthFunc func(ctx context.Context) (context.Context, error)
type LogFunc func(*requestlog.Entry)
type DescribeFunc func(chan<- *prometheus.Desc)
type CollectFunc func(chan<- prometheus.Metric)
type PromFunc func(float64)
