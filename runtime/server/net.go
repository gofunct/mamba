package server

import (
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

func (c *Service) handleGrpc(server *grpc.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			server.ServeHTTP(w, r)
		} else {
			c.r.ServeHTTP(w, r)
		}
	})
}
