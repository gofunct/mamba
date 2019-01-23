package mamba

import (
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

func (c *Command) HandleGrpc(grpcServer *grpc.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			c.Router.ServeHTTP(w, r)
		}
	})
}
