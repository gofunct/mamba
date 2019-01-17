package observe

import "github.com/prometheus/client_golang/prometheus"

type DescribeFunc func(chan<- *prometheus.Desc)
type CollectFunc func(chan<- prometheus.Metric)
type PromFunc func(float64)
