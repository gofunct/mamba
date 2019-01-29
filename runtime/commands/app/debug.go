package app

import (
	"github.com/gofunct/common/pkg/encode"
	"github.com/gofunct/common/pkg/logger/zap"
	"github.com/gofunct/mamba/runtime/config"
	mux2 "github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
)


// configCmd represents API settings command
var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "🐍  Debug current config",
	RunE: func(cmd *cobra.Command, args []string) error {
		mux := mux2.NewRouter()
		mux.HandleFunc("/debug/config", func(writer http.ResponseWriter, request *http.Request) {
			cmd.SetOutput(writer)
			cmd.Print(encode.PrettyJsonString(config.GetConfig()))
		})
		mux.HandleFunc("/debug/settings", func(writer http.ResponseWriter, request *http.Request) {
			cmd.SetOutput(writer)
			b, err := ioutil.ReadFile(V.ConfigFileUsed())
			zap.LogF("read file", err)
			cmd.Print(string(b))
		})
		mux.Handle("/debug/pprof", http.HandlerFunc(pprof.Index))
		mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		mux.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))
		zap.Debug("starting debug server", "port", "11000")
		return http.ListenAndServe(":11000", mux)
	},
}
