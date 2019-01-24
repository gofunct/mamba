package transport

import (
	"github.com/gofunct/mamba/pkg/transport/api"
	"github.com/gofunct/mamba/pkg/transport/runtime"
)

func Serve(servers ...api.Server) error {
	s := runtime.New(
		runtime.WithDefaultLogger(),
		runtime.WithServers(
			servers...,
		),
	)
	return s.Serve()
}
