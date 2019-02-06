package config

import (
	"crypto/tls"
	"github.com/gofunct/mamba/runtime/transport/api"
	"github.com/gofunct/mamba/runtime/transport/middleware"
	"net"
	"net/http"
	"os"
	"path/filepath"
	pkg_runtime "runtime"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Config contains configurations of gRPC and Gateway server.
type Config struct {
	GrpcAddr                        *Address
	GrpcInternalAddr                *Address
	GatewayAddr                     *Address
	Servers                         []api.Server
	GrpcServerUnaryInterceptors     []grpc.UnaryServerInterceptor
	GrpcServerStreamInterceptors    []grpc.StreamServerInterceptor
	GatewayServerUnaryInterceptors  []grpc.UnaryClientInterceptor
	GatewayServerStreamInterceptors []grpc.StreamClientInterceptor
	GrpcServerOption                []grpc.ServerOption
	GatewayDialOption               []grpc.DialOption
	GatewayMuxOptions               []runtime.ServeMuxOption
	GatewayServerConfig             *HTTPServerConfig
	MaxConcurrentStreams            uint32
	GatewayServerMiddlewares        []middleware.HTTPServerMiddleware
}

func CreateDefaultConfig() *Config {
	config := &Config{
		GrpcInternalAddr: &Address{
			Network: "unix",
			Addr:    "tmp/server.sock",
		},
		GatewayAddr: &Address{
			Network: "tcp",
			Addr:    ":3000",
		},
		GatewayServerConfig: &HTTPServerConfig{
			ReadTimeout:  8 * time.Second,
			WriteTimeout: 8 * time.Second,
			IdleTimeout:  2 * time.Minute,
		},
		MaxConcurrentStreams: 1000,
	}
	if pkg_runtime.GOOS == "windows" {
		config.GrpcInternalAddr = &Address{
			Network: "tcp",
			Addr:    ":5050",
		}
	}
	return config
}

// Address represents a network end point address.
type Address struct {
	Network string
	Addr    string
}

func (a *Address) CreateListener() (net.Listener, error) {
	if a.Network == "unix" {
		dir := filepath.Dir(a.Addr)
		f, err := os.Stat(dir)
		if err != nil {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return nil, errors.Wrap(err, "failed to create the directory")
			}
		} else if !f.IsDir() {
			return nil, errors.Errorf("file %q already exists", dir)
		}
	}
	lis, err := net.Listen(a.Network, a.Addr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to listen %s %s", a.Network, a.Addr)
	}
	return lis, nil
}

type HTTPServerConfig struct {
	TLSConfig         *tls.Config
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
	TLSNextProto      map[string]func(*http.Server, *tls.Conn, http.Handler)
	ConnState         func(net.Conn, http.ConnState)
}

func (c *HTTPServerConfig) ApplyTo(s *http.Server) {
	s.TLSConfig = c.TLSConfig
	s.ReadTimeout = c.ReadTimeout
	s.ReadHeaderTimeout = c.ReadHeaderTimeout
	s.WriteTimeout = c.WriteTimeout
	s.IdleTimeout = c.IdleTimeout
	s.MaxHeaderBytes = c.MaxHeaderBytes
	s.TLSNextProto = c.TLSNextProto
	s.ConnState = c.ConnState
}

func (c *Config) ServerOptions() []grpc.ServerOption {
	return append(
		[]grpc.ServerOption{
			grpc_middleware.WithUnaryServerChain(c.GrpcServerUnaryInterceptors...),
			grpc_middleware.WithStreamServerChain(c.GrpcServerStreamInterceptors...),
			grpc.MaxConcurrentStreams(c.MaxConcurrentStreams),
		},
		c.GrpcServerOption...,
	)
}

func (c *Config) ClientOptions() []grpc.DialOption {
	return append(
		[]grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithDialer(func(a string, t time.Duration) (net.Conn, error) {
				return net.Dial(c.GrpcInternalAddr.Network, a)
			}),
			grpc.WithUnaryInterceptor(
				grpc_middleware.ChainUnaryClient(c.GatewayServerUnaryInterceptors...),
			),
			grpc.WithStreamInterceptor(
				grpc_middleware.ChainStreamClient(c.GatewayServerStreamInterceptors...),
			),
		},
		c.GatewayDialOption...,
	)
}
