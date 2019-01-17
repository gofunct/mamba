package service

import (
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/google/wire"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/grpc"
	"github.com/terraform-providers/terraform-provider-kubernetes/kubernetes"
	"strings"
)

var Set = wire.NewSet(
	NewService,
	NewOptions,
)

type Service struct {
	Pattern    string
	Version    string
	Endpoint   endpoint.Endpoint
	*grpc.Server
	terraform.ResourceProvider
}

func NewService(pattern string, endpoint endpoint.Endpoint, option []grpc.ServerOption) *Service {
	s := grpc.NewServer(option...)

	return &Service{Pattern: pattern, Endpoint: endpoint, Server: s, ResourceProvider: kubernetes.Provider()}
}

// Runnable determines if the command is itself runnable.
func (c *Service) Runnable() bool {
	return c.Pattern != "" || c.Endpoint != nil || c.Server != nil || strings.Contains(c.Pattern, "/")
}

func NewOptions() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	opts = append(opts, grpc.UnaryInterceptor(kitgrpc.Interceptor))

	return opts
}
