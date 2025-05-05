package main

import (
	"flag"

	"github.com/ripple-mq/benchmark/cmd"
)

func main() {
	broker := flag.String("broker", "ripple", "Name of the broker to benchmark (e.g., ripple, kafka)")
	count := flag.Int("count", 10000, "Number of messages to produce and consume")
	messageSizeKB := flag.Int("message_size", 1, "Size of each message in KB")
	batchSize := flag.Int("batch_size", 5000, "maximum batch size in bytes")
	flag.Parse()
	cmd.Execute(*broker, *count, *messageSizeKB, *batchSize)
}
