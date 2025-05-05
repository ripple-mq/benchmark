package cmd

import (
	"github.com/ripple-mq/benchmark/tape/prometheus"
)

func PrometheusSink() {
	prometheus.ConsumerThroughput()
}
