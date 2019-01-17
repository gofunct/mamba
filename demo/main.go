package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/gofunct/mamba/app"
	"github.com/gofunct/mamba/service"
	"golang.org/x/time/rate"
	kitlog"github.com/go-kit/kit/log"
	"log"
	"os"
	"time"
)


var (
	application *app.Application
	closer func()
)

func init() {
	var err error
	application, closer, err = app.Gcp(context.Background(), "demo")
	if err != nil {
		log.Fatalf("failed to initialize application", err)
	}
}

func main() {
	defer log.Fatal(application.Execute())
	defer closer()
	// Create a single logger, which we'll use and give to other components.
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "time", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "origin", kitlog.DefaultCaller)
	}

	var sumEndpoint endpoint.Endpoint
	{
		sumEndpoint = MakeSumEndpoint(NewBasicService())
		sumEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(sumEndpoint)
		sumEndpoint = LoggingMiddleware(kitlog.With(logger, "method", "Sum"))(sumEndpoint)
	}
	svc := service.NewService("sum", sumEndpoint, nil)
	application.AddService(svc)

}
