package mamba

import (
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

func init() {
	router = mux.NewRouter()
}

var router *mux.Router

func (c *Command) handleGrpc(server *grpc.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			server.ServeHTTP(w, r)
		} else {
			router.ServeHTTP(w, r)
		}
	})
}
