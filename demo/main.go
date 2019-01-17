package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/gofunct/mamba/app"
	"golang.org/x/time/rate"
	kitlog"github.com/go-kit/kit/log"
	"log"
	"os"
	"time"
)

func main() {
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

	a, c, err := app.Gcp(context.TODO(), "demo", sumEndpoint, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer c()
	log.Fatal(a.Execute())
}
