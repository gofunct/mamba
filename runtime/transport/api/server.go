package api

import (
	"context"
	"net"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// Server is an interface for representing gRPC server implementations.
type Server interface {
	RegisterWithServer(*grpc.Server)
	RegisterWithHandler(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}

type Interface interface {
	Serve(l net.Listener) error
	Shutdown()
}
