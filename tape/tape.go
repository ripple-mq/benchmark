package tape

import (
	"encoding/json"
	"fmt"

	"github.com/ripple-mq/benchmark/tape/kafka"
	"github.com/ripple-mq/benchmark/tape/ripple"
)

func RunRipple(msgCount int, messageSizeKB int, batchSize int) {
	rippleRes := ripple.BenchmarkConsumerThroughput(msgCount, messageSizeKB, batchSize)
	PrettyPrint(rippleRes)
}

func RunKafka(msgCount int, messageSizeKB int, batchSize int) {
	kafkaRes := kafka.BenchmarkConsumerThroughput(msgCount, messageSizeKB, batchSize)
	PrettyPrint(kafkaRes)
}

func PrettyPrint(p any) {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(data))
}
