package cmd

import "github.com/ripple-mq/benchmark/tape"

func Execute(broker string, msgCount int, messageSizeKB int, batchSize int) {
	if broker == "kafka" {
		tape.RunKafka(msgCount, messageSizeKB, batchSize)
	} else if broker == "ripple" {
		tape.RunRipple(msgCount, messageSizeKB, batchSize)
	}
}
