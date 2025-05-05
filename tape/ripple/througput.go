package ripple

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ripple-mq/benchmark/pkg/utils/pen"
	"github.com/ripple-mq/go-client/wave"
)

type ConsumeThroughput struct {
	Desc               string  `json:"description"`
	Broker             string  `json:"broker"`
	ReadBatchSize      int     `json:"read_batch_size"`
	MessageSizeInBytes int     `json:"message_size_in_bytes"`
	TotalDataInBytes   float64 `json:"total_data_in_bytes"`
	ValueInMBPerSec    float64 `json:"value_in_mb_per_sec"`
}

func BenchmarkConsumerThroughput(messageCount int, messageSizeKB int, batchSize int) ConsumeThroughput {
	client := wave.NewClient[string](wave.Config{
		Topic:         uuid.NewString(),
		Bucket:        uuid.NewString(),
		Brokers:       []string{":8080"},
		ReadBatchSize: int(batchSize / (messageSizeKB * 1024)),
	})

	client.Create()
	prodCh := client.Writer()
	consCh := client.Reader()

	go func() {
		for i := 0; i < 1000100; i++ {
			prodCh <- strings.Repeat("x", 1024*messageSizeKB)
		}
	}()

	pen.SpinWheel("Producing 1M messages ", "Produced ~1M messages successfully ", 60)

	start := time.Now()

	for i := 0; i < messageCount; i++ {
		<-consCh
	}

	elapsed := time.Since(start).Seconds()
	totalMB := float64(messageCount*messageSizeKB) / 1024.0
	throughput := totalMB / elapsed

	pen.SpinWheel("Started consuming ", fmt.Sprintf("Consumed %f mb data %f", totalMB, throughput), 1)

	return ConsumeThroughput{
		Desc:               "network I/O",
		Broker:             "ripple",
		ReadBatchSize:      batchSize,
		MessageSizeInBytes: 1024 * messageSizeKB,
		TotalDataInBytes:   float64(messageCount * 1024.0),
		ValueInMBPerSec:    throughput,
	}
}
