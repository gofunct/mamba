package runtime

import (
	"database/sql"
	"fmt"
	"github.com/gofunct/mamba/runtime/health"
	"github.com/gofunct/mamba/service"
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/oklog/oklog/pkg/group"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/requestlog"
	"gocloud.dev/server"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
)

func init() {
	viper.SetConfigName("variables.hcl")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	viper.AddConfigPath(os.Getenv("PWD")+"/deploy")
	viper.AddConfigPath(os.Getenv("HOME"))
	viper.AddConfigPath(os.Getenv("."))
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	fmt.Println(viper.ConfigFileUsed())
}

var Set = wire.NewSet(
	NewService,
	trace.AlwaysSample,
	health.Set,
	service.Set,
)

type Service struct {
	db     *sql.DB
	bucket *blob.Bucket
	*server.Server
	services []*service.Service
	http.Handler
	group.Group
}

var RunGroup group.Group

func NewService(db *sql.DB, bucket *blob.Bucket, srv *server.Server, l requestlog.Logger) *Service {
	var handl http.Handler
	m := mux.NewRouter()
	m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	m.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
	handl = requestlog.NewHandler(l, m)
	return 	 &Service{db: db, bucket: bucket, Server: srv, Handler: handl, Group: RunGroup}
}

func (s *Service) Services() []*service.Service {
	return s.services
}

func (s *Service) AddService(sv *service.Service) {
	s.services = append(s.services, sv)
}

func (a *Service) ResetServices() {
	a.services = nil
}

func (s *Service) HandleGrpc(grpcServer *grpc.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			ctx, span := trace.StartSpan(r.Context(), r.URL.Host+r.URL.Path)
			defer span.End()

			r = r.WithContext(ctx)

			grpcServer.ServeHTTP(w, r)
		} else {
			ctx, span := trace.StartSpan(r.Context(), r.URL.Host+r.URL.Path)
			defer span.End()

			r = r.WithContext(ctx)
			s.Handler.ServeHTTP(w, r)
		}
	})
}

// LoadECKeyFromFile loads EC key from unencrypted PEM file.
func (s *Service) Execute() error {

	if !s.Runnable() {
		return errors.New("failed to execute- runtime is not runnable")
	}
	for _, k := range s.services {
		if k.Runnable() {
			s.Handler = s.HandleGrpc(k.Server)
			log.Printf("registered: %s", k.Pattern)
		}
	}

	log.Print("Starting server...")
	return s.ListenAndServe(":8080", s.Handler)
}

// Runnable determines if the command is itself runnable.
func (s *Service) Runnable() bool {
	return s.bucket != nil || s.db != nil || s.Handler != nil
}