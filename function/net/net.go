package net

import (
	"bufio"
	"context"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type HttpHandlerFunc func(http.ResponseWriter, *http.Request)
type RoundTripperFunc func(*http.Request) (*http.Response, error)
type ServeFunc func(addr string, h http.Handler) error
type ShutdownFunc func(ctx context.Context) error
type HijackerFunc func() (net.Conn, *bufio.ReadWriter, error)

type RouterWalkFunc func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error
type RouterBuildFunc func(map[string]string) map[string]string

type UnaryInterceptorFunc func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
type StreamInterceptorFunc func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
